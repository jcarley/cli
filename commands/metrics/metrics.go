package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	ui "github.com/gizak/termui"
)

// MetricsTransformer specifies that all concrete implementations should be
// able to transform an entire environments metrics data (group) or a single
// service metrics data (single).
type MetricsTransformer interface {
	TransformGroupCPU(*[]models.Metrics)
	TransformGroupMemory(*[]models.Metrics)
	TransformGroupNetworkIn(*[]models.Metrics)
	TransformGroupNetworkOut(*[]models.Metrics)
	TransformSingleCPU(*models.Metrics)
	TransformSingleMemory(*models.Metrics)
	TransformSingleNetworkIn(*models.Metrics)
	TransformSingleNetworkOut(*models.Metrics)
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
func CmdMetrics(svcName string, metricType MetricType, jsonFlag, csvFlag, sparkFlag, streamFlag bool, mins int, im IMetrics, is services.IServices) error {
	if streamFlag && (jsonFlag || csvFlag || mins != 1) {
		return fmt.Errorf("--stream cannot be used with a custom format and multiple records")
	}
	if mins > 1440 {
		return fmt.Errorf("--mins cannot be greater than 1440")
	}
	var mt MetricsTransformer
	if jsonFlag {
		mt = &JSONTransformer{}
	} else if csvFlag {
		buffer := &bytes.Buffer{}
		mt = &CSVTransformer{
			HeadersWritten: false,
			GroupMode:      false,
			Buffer:         buffer,
			Writer:         csv.NewWriter(buffer),
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
		//ui.UseTheme("helloworld")

		p := ui.NewPar("PRESS q TO QUIT")
		p.Border = false

		ui.Body.AddRows(
			ui.NewRow(ui.NewCol(12, 0, p)),
		)
		ui.Body.Align()
		ui.Render(ui.Body)

		mt = &SparkTransformer{
			SparkLines: map[string]*ui.Sparklines{},
		}
	} else {
		mt = &TextTransformer{}
	}
	if svcName != "" {
		service, err := is.RetrieveByLabel(svcName)
		if err != nil {
			return err
		}
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\"", svcName)
		}
		return CmdServiceMetrics(metricType, streamFlag, sparkFlag, mins, service, mt, im)
	}
	return CmdEnvironmentMetrics(metricType, streamFlag, sparkFlag, mins, mt, im)
}

func CmdEnvironmentMetrics(metricType MetricType, stream, sparkLines bool, mins int, mt MetricsTransformer, im IMetrics) error {
	done := make(chan struct{})
	go func() error {
		for {
			metrics, err := im.RetrieveEnvironmentMetrics(mins)
			if err != nil {
				done <- struct{}{}
				return err
			}
			switch metricType {
			case CPU:
				mt.TransformGroupCPU(metrics)
			case Memory:
				mt.TransformGroupMemory(metrics)
			case NetworkIn:
				mt.TransformGroupNetworkIn(metrics)
			case NetworkOut:
				mt.TransformGroupNetworkOut(metrics)
			}
			if !stream {
				break
			}
			time.Sleep(time.Minute)
		}
		done <- struct{}{}
		return nil
	}()
	if sparkLines {
		sparkLinesEventLoop()
	} else {
		<-done
	}
	return nil
}

func CmdServiceMetrics(metricType MetricType, stream, sparkLines bool, mins int, service *models.Service, mt MetricsTransformer, im IMetrics) error {
	done := make(chan struct{})
	go func() error {
		for {
			metrics, err := im.RetrieveServiceMetrics(mins, service.ID)
			if err != nil {
				done <- struct{}{}
				return err
			}
			switch metricType {
			case CPU:
				mt.TransformSingleCPU(metrics)
			case Memory:
				mt.TransformSingleMemory(metrics)
			case NetworkIn:
				mt.TransformSingleNetworkIn(metrics)
			case NetworkOut:
				mt.TransformSingleNetworkOut(metrics)
			}
			if !stream {
				break
			}
			time.Sleep(time.Minute)
		}
		done <- struct{}{}
		return nil
	}()
	if sparkLines {
		sparkLinesEventLoop()
	} else {
		<-done
	}
	return nil
}

func (m *SMetrics) RetrieveEnvironmentMetrics(mins int) (*[]models.Metrics, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/metrics?time=%dm", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, mins), headers)
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

func (m *SMetrics) RetrieveServiceMetrics(mins int, svcID string) (*models.Metrics, error) {
	headers := httpclient.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/metrics?time=%dm", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcID, mins), headers)
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
