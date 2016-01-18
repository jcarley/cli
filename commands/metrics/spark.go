package metrics

import (
	"sort"

	"github.com/catalyzeio/cli/models"
	ui "gopkg.in/gizak/termui.v1"
)

// SparkTransformer is a concrete implementation of Transformer transforming
// data using spark lines.
type SparkTransformer struct {
	Redraw     chan bool
	SparkLines map[string]*ui.Sparklines
}

func (m *SMetrics) Spark() error {
	return nil
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
