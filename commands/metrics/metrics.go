package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/commands/services"
	ui "gopkg.in/gizak/termui.v1"
)

// Transformer outlines an interface that takes in metrics data and transforms
// it into a specific data type. Suggested concrete implementations might
// include a CSVTransformer or a JSONTransformer.
type Transformer struct {
	Stream          bool
	GroupMode       bool
	Mins            int
	GroupRetriever  func(int) (*[]models.Metrics, error)
	SingleRetriever func(int) (*models.Metrics, error)
	DataTransformer MetricsTransformer

	settings *models.Settings
}

// MetricsTransformer specifies that all concrete implementations should be
// able to transform an entire environments metrics data (group) or a single
// service metrics data (single).
type MetricsTransformer interface {
	TransformGroup(*[]models.Metrics)
	TransformSingle(*models.Metrics)
}

func (transformer *Transformer) process() {
	for {
		transformer.transform()
		if !transformer.Stream {
			break
		}
		time.Sleep(time.Minute)
	}
}

func (transformer *Transformer) transform() error {
	if transformer.GroupMode {
		data, err := transformer.GroupRetriever(transformer.Mins)
		if err != nil {
			return err
		}
		transformer.DataTransformer.TransformGroup(data)
	} else {
		data, err := transformer.SingleRetriever(transformer.Mins)
		if err != nil {
			return err
		}
		transformer.DataTransformer.TransformSingle(data)
	}
	return nil
}

// go unfortunately doesn't have anything except comparisons for floats
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CmdMetrics prints out metrics for a given service or if the service is not
// specified, metrics for the entire environment are printed.
func CmdMetrics(svcName string, jsonFlag, csvFlag, sparkFlag, streamFlag bool, mins int, im IMetrics, is services.IServices) error {
	var service *models.Service
	if svcName != "" {
		service, err := is.RetrieveByLabel(svcName)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\"\n", svcName)
		}
	}
	return im.Metrics(jsonFlag, csvFlag, sparkFlag, streamFlag, mins, service)
}

func (m *SMetrics) Metrics(jsonFlag bool, csvFlag bool, sparkFlag bool, streamFlag bool, mins int, service *models.Service) error {
	if streamFlag && (jsonFlag || csvFlag || mins != 1) {
		return fmt.Errorf("--stream cannot be used with a custom format and multiple records")
	}
	var singleRetriever func(mins int) (*models.Metrics, error)
	if service != nil {
		m.Settings.ServiceID = service.ID
		singleRetriever = m.RetrieveServiceMetrics
	}
	var transformer Transformer
	redraw := make(chan bool)
	if jsonFlag {
		transformer = Transformer{
			SingleRetriever: singleRetriever,
			DataTransformer: &JSONTransformer{},
		}
	} else if csvFlag {
		buffer := &bytes.Buffer{}
		transformer = Transformer{
			SingleRetriever: singleRetriever,
			DataTransformer: &CSVTransformer{
				HeadersWritten: false,
				GroupMode:      service == nil,
				Buffer:         buffer,
				Writer:         csv.NewWriter(buffer),
			},
		}
	} else if sparkFlag {
		// the spark lines interface stays up until closed by the user, so
		// we might as well keep updating it as long as it is there
		streamFlag = true
		mins = 60
		err := ui.Init()
		if err != nil {
			return err
		}
		defer ui.Close()
		ui.UseTheme("helloworld")

		p := ui.NewPar("PRESS q TO QUIT")
		p.HasBorder = false
		p.TextFgColor = ui.Theme().SparklineTitle
		ui.Body.AddRows(
			ui.NewRow(ui.NewCol(12, 0, p)),
		)

		transformer = Transformer{
			SingleRetriever: singleRetriever,
			DataTransformer: &SparkTransformer{
				Redraw:     redraw,
				SparkLines: make(map[string]*ui.Sparklines),
			},
		}
	} else {
		transformer = Transformer{
			SingleRetriever: singleRetriever,
			DataTransformer: &TextTransformer{},
		}
	}
	transformer.GroupRetriever = m.RetrieveEnvironmentMetrics
	transformer.Stream = streamFlag
	transformer.GroupMode = service == nil
	transformer.Mins = mins
	transformer.settings = m.Settings

	// TODO why is this here? -> helpers.SignIn(m.Settings)

	if sparkFlag {
		go transformer.process()

		ui.Body.Align()
		ui.Render(ui.Body)

		quit := make(chan bool)
		go maintainSparkLines(redraw, quit)
		<-quit
	} else {
		transformer.process()
	}
	return nil
}

func (m *SMetrics) RetrieveEnvironmentMetrics(mins int) (*[]models.Metrics, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/metrics?time=%d", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, mins), headers)
	if err != nil {
		return nil, err
	}
	var metrics []models.Metrics
	err = httpclient.ConvertResp(resp, statusCode, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

func (m *SMetrics) RetrieveServiceMetrics(mins int) (*models.Metrics, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/metrics?time=%d", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, m.Settings.ServiceID, mins), headers)
	if err != nil {
		return nil, err
	}
	var metrics models.Metrics
	err = httpclient.ConvertResp(resp, statusCode, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}
