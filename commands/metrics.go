package commands

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
	ui "github.com/gizak/termui"
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

// TextTransformer is a concrete implementation of Transformer transforming
// data into plain text.
type TextTransformer struct{}

// CSVTransformer is a concrete implementation of Transformer transforming
// data into CSV.
type CSVTransformer struct {
	HeadersWritten bool
	GroupMode      bool // I cant figure out a good way to access the Transformer.GroupMode from a CSVTransformer instance, so I copied it
	Buffer         *bytes.Buffer
	Writer         *csv.Writer
}

// JSONTransformer is a concrete implementation of Transformer transforming
// data into JSON.
type JSONTransformer struct{}

// SparkTransformer is a concrete implementation of Transformer transforming
// data using spark lines.
type SparkTransformer struct {
	Redraw     chan bool
	SparkLines map[string]*ui.Sparklines
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

// TransformGroup transforms an entire environment's metrics data into text
// format.
func (text *TextTransformer) TransformGroup(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		fmt.Printf("%s:\n", metric.ServiceName)
		text.TransformSingle(&metric)
	}
}

// TransformSingle transforms a single service's metrics data into text
// format.
func (text *TextTransformer) TransformSingle(metric *models.Metrics) {
	prefix := "    "
	for _, job := range *metric.Jobs {
		for _, data := range *job.MetricsData {
			ts := time.Unix(data.TS, 0)
			fmt.Printf("%s%s | %8s (%s) | CPU: %6.2fs (%5.2f%%) | Net: RX: %.2f KB TX: %.2f KB | Mem: %.2f KB | Disk: %.2f KB read / %.2f KB write\n",
				prefix,
				fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
				job.Type,
				job.ID,
				float64(data.CPU.Usage)/1000000000.0,
				float64(data.CPU.Usage)/1000000000.0/60.0*100.0,
				math.Ceil(data.Network.RXKb),
				math.Ceil(data.Network.TXKb),
				math.Ceil(data.Memory.Avg/1024.0),
				math.Ceil(data.DiskIO.Read/1024.0),
				math.Ceil(data.DiskIO.Write/1024.0))
		}
	}
}

// WriteHeaders prints out the csv headers to the console.
func (csv *CSVTransformer) WriteHeaders() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "type", "job_id", "cpu_usage", "cpu_percentage", "rx_bytes", "tx_bytes", "memory", "disk_read", "disk_write"}
		if csv.GroupMode {
			headers = append([]string{"service_label", "service_id"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

// TransformGroup transforms an entire environment's metrics data into csv
// format.
func (csv *CSVTransformer) TransformGroup(metrics *[]models.Metrics) {
	csv.WriteHeaders()
	for _, metric := range *metrics {
		csv.TransformSingle(&metric)
	}
	csv.Writer.Flush()
	fmt.Println(csv.Buffer.String())
}

// TransformSingle transforms a single service's metrics data into csv
// format.
func (csv *CSVTransformer) TransformSingle(metric *models.Metrics) {
	csv.WriteHeaders()
	for _, job := range *metric.Jobs {
		for _, data := range *job.MetricsData {
			row := []string{
				strconv.FormatInt(data.TS, 10),
				job.Type,
				job.ID,
				strconv.FormatFloat(math.Ceil(data.CPU.Usage/1000000000.0), 'f', -1, 64),
				strconv.FormatFloat(math.Ceil(data.CPU.Usage/1000000000.0/60.0*100.0), 'f', -1, 64),
				strconv.FormatFloat(math.Ceil(data.Network.RXKb), 'f', -1, 64),
				strconv.FormatFloat(math.Ceil(data.Network.TXKb), 'f', -1, 64),
				strconv.FormatFloat(math.Ceil(data.Memory.Avg/1024.0), 'f', -1, 64),
				strconv.FormatFloat(math.Ceil(data.DiskIO.Read/1024.0), 'f', -1, 64),
				strconv.FormatFloat(math.Ceil(data.DiskIO.Write/1024.0), 'f', -1, 64),
			}
			if csv.GroupMode {
				row = append([]string{metric.ServiceName, metric.ServiceID}, row...)
			}
			csv.Writer.Write(row)
		}
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		fmt.Println(csv.Buffer.String())
	}
}

// TransformGroup transforms an entire environment's metrics data into json
// format.
func (j *JSONTransformer) TransformGroup(metrics *[]models.Metrics) {
	b, _ := json.MarshalIndent(metrics, "", "    ")
	fmt.Println(string(b))
}

// TransformSingle transforms a single service's metrics data into json
// format.
func (j *JSONTransformer) TransformSingle(metric *models.Metrics) {
	b, _ := json.MarshalIndent(metric, "", "    ")
	fmt.Println(string(b))
}

// go unfortunately doesn't have anything except comparisons for floats
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// deltas computes the change between each datapoint. Since there will be
// n-1 deltas, add a zero to the front. While not perfectly accurate, the
// leading zero sets the minimum to always be zero.
func (spark *SparkTransformer) deltas(items []int) []int {
	var newItems []int
	newItems = append(newItems, 0)
	for i := 0; i < len(items)-1; i++ {
		newItems = append(newItems, max(0, items[i+1]-items[i]))
	}
	return newItems
}

// TransformGroup transforms an entire environment's metrics data using
// spark lines.
func (spark *SparkTransformer) TransformGroup(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		spark.TransformSingle(&metric)
	}
}

// TransformSingle transforms a single service's metrics data using spark lines.
func (spark *SparkTransformer) TransformSingle(metric *models.Metrics) {
	for _, job := range *metric.Jobs {
		sparkData := make(map[string][]int)
		for i := len(*job.MetricsData) - 1; i >= 0; i-- {
			data := (*job.MetricsData)[i]
			sparkData["CPU"] = append(sparkData["CPU"], int(data.CPU.Usage))
			sparkData["Disk Read"] = append(sparkData["Disk Read"], int(data.DiskIO.Read))
			sparkData["Disk Write"] = append(sparkData["Disk Write"], int(data.DiskIO.Write))
			sparkData["Memory"] = append(sparkData["Memory"], int(data.Memory.Avg))
			sparkData["Net RX"] = append(sparkData["Net RX"], int(data.Network.RXKb))
			sparkData["Net TX"] = append(sparkData["Net TX"], int(data.Network.TXKb))
		}
		sparkData["Net RX"] = spark.deltas(sparkData["Net RX"])
		sparkData["Net TX"] = spark.deltas(sparkData["Net TX"])
		var sortedKeys []string
		for k := range sparkData {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
		var slData [][]int
		for _, key := range sortedKeys {
			value := sparkData[key]
			slData = append(slData, value)
		}
		var sparkLines = spark.SparkLines[metric.ServiceName]
		if sparkLines == nil {
			sparkLines = addSparkLine(metric.ServiceName, slData)
			spark.SparkLines[metric.ServiceName] = sparkLines
		} else {
			for i := range sparkLines.Lines {
				sparkLines.Lines[i].Data = slData[i]
			}
		}
		spark.Redraw <- true
	}
}

// Metrics prints out metrics for a given service or if the service is not
// specified, metrics for the entire environment are printed.
func Metrics(serviceLabel string, jsonFlag bool, csvFlag bool, sparkFlag bool, streamFlag bool, mins int, settings *models.Settings) {
	if streamFlag && (jsonFlag || csvFlag || mins != 1) {
		fmt.Println("--stream cannot be used with a custom format and multiple records")
		os.Exit(1)
	}
	var singleRetriever func(mins int, settings *models.Settings) *models.Metrics
	if serviceLabel != "" {
		service := helpers.RetrieveServiceByLabel(serviceLabel, settings)
		if service == nil {
			fmt.Printf("Could not find a service with the label \"%s\"\n", serviceLabel)
			os.Exit(1)
		}
		settings.ServiceID = service.ID
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
				GroupMode:      serviceLabel == "",
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
	transformer.GroupMode = serviceLabel == ""
	transformer.Mins = mins
	transformer.settings = settings

	helpers.SignIn(settings)

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
}

func addSparkLine(serviceName string, data [][]int) *ui.Sparklines {
	titles := []string{"CPU", "Disk Read", "Disk Write", "Memory", "Net RX", "Net TX"}
	var sparkLines []ui.Sparkline
	for i, title := range titles {
		sparkLine := ui.NewSparkline()
		sparkLine.Height = 1
		sparkLine.Data = data[i]
		sparkLine.Title = title
		sparkLines = append(sparkLines, sparkLine)
	}
	sp := ui.NewSparklines(sparkLines...)
	sp.Height = 14
	sp.Border.Label = serviceName

	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(12, 0, sp)),
	)

	ui.Body.Align()
	ui.Render(sp)

	return sp
}

func maintainSparkLines(redraw chan bool, quit chan bool) {
	evt := ui.EventCh()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				quit <- true
				return
			}
			if e.Type == ui.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
				go func() { redraw <- true }()
			}
		case <-redraw:
			ui.Render(ui.Body)
		}
	}
}
