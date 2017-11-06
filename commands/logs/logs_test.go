package logs

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/commands/services"
	"github.com/daticahealth/cli/commands/sites"
	"github.com/daticahealth/cli/config"
	"github.com/daticahealth/cli/lib/jobs"
	"github.com/daticahealth/cli/models"
	"github.com/daticahealth/cli/test"
)

type SLogsMock struct {
	Settings *models.Settings
}

func (l *SLogsMock) RetrieveElasticsearchVersion(domain string) (string, error) {
	return "5", nil
}

func (l *SLogsMock) Output(queryString, domain string, generator queryGenerator, from int, startTimestamp, endTimestamp time.Time, hostNames []string, fileName string) (int, error) {
	appLogsIdentifier := "source"
	appLogsValue := "app"
	if strings.HasPrefix(domain, "csb01") {
		appLogsIdentifier = "syslog_program"
		appLogsValue = "supervisord"
	}

	logrus.Println("        @timestamp       -        message")
	for {
		queryBytes, err := generator(queryString, appLogsIdentifier, appLogsValue, startTimestamp, from, hostNames, fileName)
		if err != nil {
			return -1, fmt.Errorf("Error generating query: %s", err)
		} else if queryBytes == nil || len(queryBytes) == 0 {
			return -1, errors.New("Error generating query")
		}

		var logs models.Logs
		var hits models.Hits
		logHits := []models.LogHits{}
		for i := 0; i < 3; i++ {
			var logHit models.LogHits
			logHit.ID = fmt.Sprintf("log_%d", i)
			logHit.Score = 2.3
			logHit.Index = "1"
			logHit.Type = "THeHittenistLogHit"
			logHit.Fields = make(map[string][]string)
			logHit.Fields["@timestamp"] = []string{fmt.Sprintf("2017-10-11T15:04:0%d.999999999Z07:00", i)}
			logHit.Fields["message"] = []string{fmt.Sprintf("Wow so log %d", i)}
			logHits = append(logHits, logHit)
		}
		hits.Hits = &logHits
		hits.MaxScore = 2.3
		hits.Total = 1
		logs.Hits = &hits

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

func (l *SLogsMock) Stream(queryString, domain string, generator queryGenerator, from int, timestamp time.Time, hostNames []string, fileName string) error {
	//Don't want to run stream forever in test
	for i := 0; i < 2; i++ {
		f, err := l.Output(queryString, domain, generator, from, timestamp, time.Now(), hostNames, fileName)
		if err != nil {
			return err
		}
		from = f
		time.Sleep(config.LogPollTime * time.Second)
	}
	return nil
}

func (l *SLogsMock) Watch(queryString, domain string) error {
	//TODO: Mock it better?
	return errors.New("Run Stream")
}

func muxSetup(mux *http.ServeMux, t *testing.T, serviceType string, createdAt []string, query *CMDLogQuery) {
	mux.HandleFunc("/environments/"+test.EnvID+"/services/",
		// Retrieve services
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"%s","name":"%s","type":"%s","redeployable":false},{"id":"%s","label":"service_proxy","name":"service proxy","redeployable":true}]`, test.SvcID, test.SvcLabel, serviceType, serviceType, test.SvcIDAlt))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/jobs/"+query.JobID,
		//Retrieve job
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			if query.JobID == test.JobID {
				jobJSON := fmt.Sprintf(`{"id":"%s","type":"%s","target":"%s","status":"happy", "created_at":"%s"}`, test.JobID, "deploy", test.Target, createdAt[0])
				fmt.Fprint(w, jobJSON)
			} else {
				fmt.Fprint(w, "")
			}
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/jobs",
		// RetrieveByTarget/Type
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			var jobs []string
			for i, created := range createdAt {
				jobID := test.JobID
				if i > 0 {
					jobID = test.JobIDAlt
				}
				jobs = append(jobs, fmt.Sprintf(`{"id":"%s","type":"%s","target":"%s","status":"happy","created_at":"%s"}`, jobID, "worker", test.Target, created))
			}
			jobsJSON := fmt.Sprintf("[%s]", strings.Join(jobs, ","))
			fmt.Fprint(w, jobsJSON)
		},
	)
	mux.HandleFunc("/environments",
		// Retrieve environment list
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			if r.Header.Get("X-Pod-ID") == test.Pod {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
			} else {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvIDAlt, test.EnvNameAlt, test.NamespaceAlt, test.OrgIDAlt))
			}
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID,
		// Retrieve environment by ID
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcIDAlt+"/sites",
		// Retrieve site by ID
		func(w http.ResponseWriter, r *http.Request) {
			test.AssertEquals(t, r.Method, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":%d, "name":"%s"}]`, 123, test.Namespace+".supersite"))
		},
	)
}

func TestLogsService(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestLogsJobID(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		JobID:   test.JobID,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestLogsTarget(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		Target:  test.Target,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate, test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestLogsStream(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  true,
		Service: test.SvcLabel,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestLogsBadService(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: "BadServiceLabel",
	}
	muxSetup(mux, t, "code", []string{test.BadDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := fmt.Sprintf("Cannot find the specified service \"%s\".", cmdQuery.Service)
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsBadRequestMissingService(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:  "",
		Follow: false,
		Target: test.Target,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := "You must specify a code service to query the logs for a particular target"
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsBadRequestTargetAndJobID(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:  "",
		Follow: false,
		Target: test.Target,
		JobID:  test.JobID,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := "Specifying \"--job-id\" in combination with \"--target\" is unsupported."
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsBadRequestTargetNonCodeService(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		Target:  test.Target,
	}
	muxSetup(mux, t, "database", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := "Cannot specifiy a target for a non-code service type"
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsBadServiceWithTarget(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: "BadServiceLabel",
		Target:  test.Target,
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := fmt.Sprintf("Cannot find the specified service \"%s\".", cmdQuery.Service)
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsBadJobID(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		JobID:   "BadJobID",
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := fmt.Sprintf("Cannot find the specified job \"%s\".", cmdQuery.JobID)
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsBadTarget(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		Target:  "BadTarget",
	}
	muxSetup(mux, t, "code", []string{test.GoodDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := fmt.Sprintf("Cannot find any jobs with target \"%s\" for service \"%s\"", cmdQuery.Target, test.SvcID)
	if err == nil {
		t.Fatalf("Expected: %s\n", expectedErr)
	} else if err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsNoValidHostNames(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		Target:  test.Target,
	}
	muxSetup(mux, t, "code", []string{test.BadDate, test.BadDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}
	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	expectedErr := fmt.Sprintf(`All %d jobs for the service "%s"%s do not have valid hostnames to allow their logs to be queried. Redeploy the service if you would like to use this functionality.`, 2, cmdQuery.Service, fmt.Sprintf(` that have a target of "%s"`, cmdQuery.Target))
	if err == nil || err.Error() != expectedErr {
		t.Fatalf("Expected: %s\nGot: %s", expectedErr, err)
	}
}

func TestLogsOneValidHostName(t *testing.T) {
	mux, server, baseURL := test.Setup()
	defer test.Teardown(server)
	settings := test.GetSettings(baseURL.String())
	cmdQuery := CMDLogQuery{
		Query:   "",
		Follow:  false,
		Service: test.SvcLabel,
		Target:  test.Target,
	}
	muxSetup(mux, t, "code", []string{test.BadDate, test.GoodDate, test.BadDate, test.BadDate}, &cmdQuery)

	ilogs := &SLogsMock{
		Settings: settings,
	}

	err := CmdLogs(&cmdQuery, settings.EnvironmentID, settings, ilogs, &test.FakePrompts{}, environments.New(settings), services.New(settings), jobs.New(settings), sites.New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
