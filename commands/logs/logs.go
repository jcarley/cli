package logs

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/sites"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
)

const size = 50

type esVersion struct {
	Number string `json:"number"`
}

type CMDLogQuery struct {
	Query   string // default *
	Follow  bool
	Hours   int
	Minutes int
	Seconds int
	Service string
	JobID   string
	Target  string
}

// CmdLogs is a way to stream logs from Kibana to your local terminal. This is
// useful because Kibana is hard to look at because it splits every single
// log statement into a separate block that spans multiple lines so it's
// not very cohesive. This is intended to be similar to the `heroku logs`
// command.
func CmdLogs(query *CMDLogQuery, envID string, settings *models.Settings, il ILogs, ip prompts.IPrompts, ie environments.IEnvironments, is services.IServices, ij jobs.IJobs, isites sites.ISites) error {
	if query.Follow && (query.Hours > 0 || query.Minutes > 0 || query.Seconds > 0) {
		return fmt.Errorf("Specifying \"-f\" in combination with \"--hours\", \"--minutes\", or \"--seconds\" is unsupported.")
	}
	if len(query.JobID) > 0 && len(query.Target) > 0 {
		fmt.Errof("Specifying \"--job-id\" in combination with \"--target\" is unsupported.")
	}
	if len(query.JobID) > 0 && len(query.Service) == 0 {
		return fmt.Errorf("You must specify a service to query the logs for a particular job.")
	}
	if len(query.Target) > 0 && len(query.Service) == 0 {
		return fmt.Errorf("You must specify a least a service to query the logs for a particular target")
	}
	var svc *models.Service
	var hostNames []string
	var fileName string
	if len(query.JobID) > 0 || len(query.Target) > 0 {
		t := time.Unix(startOfReadableHostNames, 0)
		var err error
		svc, err = is.RetrieveByLabel(query.Service)
		if err != nil {
			return err
		}
		if svc == nil {
			return fmt.Erroff("Cannot find the specified service \"%s\".", query.Service)
		}
		if len(query.JobID) > 0 {
			job, err := ij.Retrieve(query.JobID, svc.ID, true)
			if err != nil {
				return err
			}
			if job == nil {
				return fmt.Erroff("Cannot find the specified job \"%s\".", query.JobID)
			}
			if len(job.Spec.Description.HostName) == 0 {
				return fmt.Errorf("This job, \"%s\", does not have a valid host name and therefore its logs are not marked with its \"ID\".")
			}
			hostNames = []string{job.Spec.Description.HostName}
		} else if len(query.Target) > 0 {
			jobs, err := ij.RetrieveByTarget(svc.ID, query.Target, 1, 25)
			if err != nil {
				return err
			}
			if jobs == nil {
				return fmt.Errorf("Cannot find any jobs with target \"%s\" for service \"%s\"", query.Target, svc.ID)
			}
			var totalCount, badCount int
			for _, j := range jobs {
				hostName := j.Spec.Description.HostName
				if len(hostName) > 0 {
					hostNames := append(hostNames, hostName)
				} else {
					badCount++
				}
				totalCount++
			}
			if badCount > 0 {
				err := prompts.YesNo("", fmt.Sprintf("Of the %d jobs for the service \"%s\" that have a target of \"%s\" %d do not have a valid hostname to allow their logs to be queried. Would you like to proceed anyways? (y/n)", totalCount, query.Service, query.Target, badCount))
				if err != nil {
					return err
				}
			}
		}
	} else if len(query.Service) {

	}
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	sites, err := isites.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	domain := ""
	for _, site := range *sites {
		if strings.HasPrefix(site.Name, env.Namespace) {
			domain = site.Name
			break
		}
	}
	if domain == "" {
		return errors.New("Could not determine the fully qualified domain name of your environment. Please contact Datica Support at https://datica.com/support with this error message to resolve this issue.")
	}
	version, err := il.RetrieveElasticsearchVersion(domain)
	if err != nil {
		version = ""
	}
	generator := chooseQueryGenerator(version)
	if follow {
		if err = il.Watch(queryString, domain); err != nil {
			logrus.Debugf("Error attempting to stream logs from logwatch: %s", err)
		} else {
			return nil
		}
	}
	from := 0
	offset := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	timestamp := time.Now().In(time.UTC).Add(-1 * offset)
	from, err = il.Output(queryString, domain, generator, from, timestamp, time.Now())
	if err != nil {
		return err
	}
	if follow {
		return il.Stream(queryString, domain, generator, from, timestamp)
	}
	return nil
}

func (l *SLogs) RetrieveElasticsearchVersion(domain string) (string, error) {
	headers := map[string][]string{"Cookie": {"sessionToken=" + url.QueryEscape(l.Settings.SessionToken)}}
	resp, statusCode, err := l.Settings.HTTPManager.Get(nil, fmt.Sprintf("https://%s/__es/", domain), headers)
	if err != nil {
		return "", err
	}
	var wrapper struct {
		Version esVersion `json:"version"`
	}
	err = l.Settings.HTTPManager.ConvertResp(resp, statusCode, &wrapper)
	if err != nil {
		return "", err
	}
	return wrapper.Version.Number, nil
}

func (l *SLogs) Output(queryString, domain string, generator queryGenerator, from int, startTimestamp, endTimestamp time.Time) (int, error) {
	appLogsIdentifier := "source"
	appLogsValue := "app"
	if strings.HasPrefix(domain, "csb01") {
		appLogsIdentifier = "syslog_program"
		appLogsValue = "supervisord"
	}

	urlString := fmt.Sprintf("https://%s/__es", domain)

	headers := map[string][]string{"Cookie": {"sessionToken=" + url.QueryEscape(l.Settings.SessionToken)}}

	logrus.Println("        @timestamp       -        message")
	for {
		queryBytes := generator(queryString, appLogsIdentifier, appLogsValue, startTimestamp, from)

		resp, statusCode, err := l.Settings.HTTPManager.Get(queryBytes, fmt.Sprintf("%s/_search", urlString), headers)
		if err != nil {
			return from, err
		}
		var logs models.Logs
		err = l.Settings.HTTPManager.ConvertResp(resp, statusCode, &logs)
		if err != nil {
			return from, err
		}

		end := time.Time{}
		for _, lh := range *logs.Hits.Hits {
			logrus.Printf("%s - %s", lh.Fields["@timestamp"][0], lh.Fields["message"][0])
			end, _ = time.Parse(time.RFC3339Nano, lh.Fields["@timestamp"][0])
		}
		amount := len(*logs.Hits.Hits)

		from += len(*logs.Hits.Hits)
		// TODO this will infinite loop if it always retrieves `size` hits
		// and it fails to parse the end timestamp. very small window of opportunity.
		if amount < size || end.After(endTimestamp) {
			break
		}
		time.Sleep(config.JobPollTime * time.Second)
	}
	return from, nil
}

func (l *SLogs) Stream(queryString, domain string, generator queryGenerator, from int, timestamp time.Time) error {
	for {
		f, err := l.Output(queryString, domain, generator, from, timestamp, time.Now())
		if err != nil {
			return err
		}
		from = f
		time.Sleep(config.LogPollTime * time.Second)
	}
}
