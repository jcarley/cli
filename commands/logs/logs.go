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
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
)

const size = 50

type esVersion struct {
	Number string `json:"number"`
}

// CmdLogs is a way to stream logs from Kibana to your local terminal. This is
// useful because Kibana is hard to look at because it splits every single
// log statement into a separate block that spans multiple lines so it's
// not very cohesive. This is intended to be similar to the `heroku logs`
// command.
func CmdLogs(queryString string, follow bool, hours, minutes, seconds int, envID string, settings *models.Settings, il ILogs, ip prompts.IPrompts, ie environments.IEnvironments, is services.IServices, isites sites.ISites) error {
	if follow && (hours > 0 || minutes > 0 || seconds > 0) {
		return fmt.Errorf("Specifying \"-f\" in combination with \"--hours\", \"--minutes\", or \"--seconds\" is unsupported.")
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
