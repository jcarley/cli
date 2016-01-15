package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
	ui "gopkg.in/gizak/termui.v1"
)

// Transformer outlines an interface that takes in metrics data and transforms
// it into a specific data type. Suggested concrete implementations might
// include a CSVTransformer or a JSONTransformer.
type Transformer struct {
	Stream          bool
	GroupMode       bool
	Mins            int
	GroupRetriever  func(int, *models.Settings) *[]models.Metrics
	SingleRetriever func(int, *models.Settings) *models.Metrics
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

func (transformer *Transformer) transform() {
	if transformer.GroupMode {
		transformer.DataTransformer.TransformGroup(transformer.GroupRetriever(transformer.Mins, transformer.settings))
	} else {
		transformer.DataTransformer.TransformSingle(transformer.SingleRetriever(transformer.Mins, transformer.settings))
	}
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
func CmdMetrics(svcName string, jsonFlag, csvFlag, sparkFlag, streamFlag bool, mins int, im IMetrics) error {
	return im.Metrics(svcName, jsonFlag, csvFlag, sparkFlag, streamFlag, mins, im)
}

func (m *SMetrics) Metrics(svcName string, jsonFlag bool, csvFlag bool, sparkFlag bool, streamFlag bool, mins int, im IMetrics) error {
	if streamFlag && (jsonFlag || csvFlag || mins != 1) {
		return fmt.Errorf("--stream cannot be used with a custom format and multiple records")
	}
	var singleRetriever func(mins int, settings *models.Settings) *models.Metrics
	if svcName != "" {
		service := helpers.RetrieveServiceByLabel(svcName, m.Settings)
		if service == nil {
			return fmt.Errorf("Could not find a service with the label \"%s\"\n", svcName)
		}
		m.Settings.ServiceID = service.ID
		singleRetriever = helpers.RetrieveServiceMetrics
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
				GroupMode:      svcName == "",
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
			fmt.Println(err.Error())
			os.Exit(1)
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
	transformer.GroupRetriever = helpers.RetrieveEnvironmentMetrics
	transformer.Stream = streamFlag
	transformer.GroupMode = svcName == ""
	transformer.Mins = mins
	transformer.settings = m.Settings

	helpers.SignIn(m.Settings)

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
