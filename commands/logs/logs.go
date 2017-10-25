package logs

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
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

// type CMDLogQuery struct {
// 	Query    string // default *
// 	Follow   bool
// 	Hours    int
// 	Minutes  int
// 	Seconds  int
// 	Service  string
// 	JobID    string
// 	Target   string
// 	FileName string
// }

// CmdLogs is a way to stream logs from Kibana to your local terminal. This is
// useful because Kibana is hard to look at because it splits every single
// log statement into a separate block that spans multiple lines so it's
// not very cohesive. This is intended to be similar to the `heroku logs`
// command.
//TODO: Take out envID? Already present in settings
func CmdLogs(query *CMDLogQuery, envID string, settings *models.Settings, il ILogs, ip prompts.IPrompts, ie environments.IEnvironments, is services.IServices, ij jobs.IJobs, isites sites.ISites) error {
	if query.Follow && (query.Hours > 0 || query.Minutes > 0 || query.Seconds > 0) {
		return fmt.Errorf("Specifying \"-f\" in combination with \"--hours\", \"--minutes\", or \"--seconds\" is unsupported.")
	}
	if len(query.JobID) > 0 && len(query.Target) > 0 {
		return fmt.Errorf("Specifying \"--job-id\" in combination with \"--target\" is unsupported.")
	}
	if len(query.JobID) > 0 && len(query.Service) == 0 {
		return fmt.Errorf("You must specify a service to query the logs for a particular job.")
	}
	if len(query.Target) > 0 && len(query.Service) == 0 {
		return fmt.Errorf("You must specify a code service to query the logs for a particular target")
	}
	var svcID string
	var hostNames []string
	//TODO: Get from service label if no jobID or target defined
	//Replace all non alphanumeric characters with underscores
	//Must check service type, not all filnames are defined the same
	//If not code or  custom service, then fall back to hostName
	//Get all deploy jobs for service, check if host names exist. If not, then can't do this command
	var fileName string
	if len(query.JobID) > 0 || len(query.Target) > 0 {
		// t := time.Unix(startOfReadableHostNames, 0)
		svc, err := is.RetrieveByLabel(query.Service)
		if err != nil {
			return err
		}
		if svc == nil {
			return fmt.Errorf("Cannot find the specified service \"%s\".", query.Service)
		}
		svcID = svc.ID
		if len(query.JobID) > 0 {
			job, err := ij.Retrieve(query.JobID, svc.ID, true)
			if err != nil {
				return err
			}
			if job == nil || job.ID != query.JobID {
				return fmt.Errorf("Cannot find the specified job \"%s\".", query.JobID)
			}
			if len(job.Spec.Description.HostName) == 0 {
				return fmt.Errorf("This job, \"%s\", does not have a valid host name and therefore its logs are not marked with its \"ID\".", job.ID)
			}
			hostNames = []string{job.Spec.Description.HostName}
		} else if len(query.Target) > 0 {
			if svc.Type != "code" {
				return fmt.Errorf("Cannot specifiy a target for a non-code service type")
			}
			jobs, err := ij.RetrieveByTarget(svcID, query.Target, 1, 25)
			if err != nil {
				return err
			}
			if jobs == nil || len(*jobs) == 0 {
				return fmt.Errorf("Cannot find any jobs with target \"%s\" for service \"%s\"", query.Target, svcID)
			}
			hostNames, err = getHostNames(jobs, query, ip)
			if err != nil {
				return err
			}
		}
	} else if len(query.Service) > 0 {
		svc, err := is.RetrieveByLabel(query.Service)
		if err != nil {
			return err
		}
		if svc == nil {
			return fmt.Errorf("Cannot find the specified service \"%s\".", query.Service)
		}
		svcID = svc.ID
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return err
		}
		switch svc.Type {
		case "app":
			{
				fileName = "/data/log/app/" + reg.ReplaceAllString(svc.Label, "_") + "/current"
			}
		case "custom":
			{
				fileName = "/data/log/custom/" + reg.ReplaceAllString(svc.Label, "_") + "/current"
			}
		default:
			{
				jobs, err := ij.RetrieveByType(svcID, "deploy", 1, 25)
				if err != nil {
					return err
				}
				if jobs == nil {
					return fmt.Errorf("Cannot find any deploy jobs for service \"%s\"", svcID)
				}
				hostNames, err = getHostNames(jobs, query, ip)
				if err != nil {
					return err
				}
			}
		}
	}
	env, err := ie.Retrieve(envID)
	if err != nil {
		fmt.Println("Err retruevug env")
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
	if query.Follow {
		if err = il.Watch(query.Query, domain); err != nil {
			logrus.Debugf("Error attempting to stream logs from logwatch: %s", err)
		} else {
			return nil
		}
	}
	from := 0
	offset := time.Duration(query.Hours)*time.Hour + time.Duration(query.Minutes)*time.Minute + time.Duration(query.Seconds)*time.Second
	timestamp := time.Now().In(time.UTC).Add(-1 * offset)
	from, err = il.Output(query.Query, domain, generator, from, timestamp, time.Now(), hostNames, fileName)
	if err != nil {
		return err
	}
	if query.Follow {
		return il.Stream(query.Query, domain, generator, from, timestamp, hostNames, fileName)
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

func (l *SLogs) Output(queryString, domain string, generator queryGenerator, from int, startTimestamp, endTimestamp time.Time, hostNames []string, fileName string) (int, error) {
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
		queryBytes, err := generator(queryString, appLogsIdentifier, appLogsValue, startTimestamp, from, hostNames, fileName)
		if err != nil {
			return -1, fmt.Errorf("Error generating query: %s", err)
		} else if queryBytes == nil || len(queryBytes) == 0 {
			return -1, errors.New("Error generating query")
		}

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
			if len(lh.Fields) != 0 {
				logrus.Printf("%s - %s", lh.Fields["@timestamp"][0], lh.Fields["message"][0])
				end, _ = time.Parse(time.RFC3339Nano, lh.Fields["@timestamp"][0])
			}
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

func (l *SLogs) Stream(queryString, domain string, generator queryGenerator, from int, timestamp time.Time, hostNames []string, fileName string) error {
	for {
		f, err := l.Output(queryString, domain, generator, from, timestamp, time.Now(), hostNames, fileName)
		if err != nil {
			return err
		}
		from = f
		time.Sleep(config.LogPollTime * time.Second)
	}
}

func getHostNames(jobs *[]models.Job, query *CMDLogQuery, ip prompts.IPrompts) ([]string, error) {
	var totalCount, badCount int
	var hostNames []string
	fmt.Printf("Jobs: %d\n", len(*jobs))
	for _, j := range *jobs {
		// TODO: Spec not getting filled
		fmt.Printf("SPEC: %s\n", j.Spec)
		fmt.Printf("Job: %v\n", j)
		if j.Spec != nil {
			hostName := j.Spec.Description.HostName
			if len(hostName) > 0 {
				hostNames = append(hostNames, hostName)
			} else {
				badCount++
			}
		} else {
			badCount++
		}
		totalCount++
	}
	if badCount > 0 {
		targetString := ""
		if len(query.Target) > 0 {
			targetString = fmt.Sprintf(` that have a target of "%s"`, query.Target)
		}
		if badCount == totalCount {
			return nil, fmt.Errorf(`All %d jobs for the service "%s"%s do not have valid hostnames to allow their logs to be queried. Unable to proceed`, totalCount, query.Service, targetString)
		}
		prompt := fmt.Sprintf("Of the %d jobs for the service \"%s\"%s %d do not have a valid hostname to allow their logs to be queried. Would you like to proceed anyways?", totalCount, query.Service, targetString, badCount)
		err := ip.YesNo("(y/n)", prompt)
		if err != nil {
			return nil, err
		}
	}
	return hostNames, nil
}
