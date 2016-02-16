package metrics

import (
	"github.com/catalyzeio/cli/models"
	ui "github.com/gizak/termui"
)

// SparkTransformer is a concrete implementation of Transformer transforming
// data using spark lines.
type SparkTransformer struct {
	SparkLines map[string]*ui.Sparklines
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
func (spark *SparkTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		spark.TransformSingleCPU(&metric)
	}
}

func (spark *SparkTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		spark.TransformSingleMemory(&metric)
	}
}

func (spark *SparkTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		spark.TransformSingleNetworkIn(&metric)
	}
}

func (spark *SparkTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		spark.TransformSingleNetworkOut(&metric)
	}
}

// TransformSingle transforms a single service's metrics data using spark lines.
func (spark *SparkTransformer) TransformSingleCPU(metric *models.Metrics) {
	var cpuMin []int
	var cpuMax []int
	var cpuAvg []int
	var cpuTotal []int
	for _, data := range *metric.Data.CPULoad {
		cpuMin = append(cpuMin, int(data.Min/1000.0))
		cpuMax = append(cpuMax, int(data.Max/1000.0))
		cpuAvg = append(cpuAvg, int(data.AVG/1000.0))
		cpuTotal = append(cpuTotal, int(data.Total/1000.0))
	}
	var sparkLines = spark.SparkLines[metric.ServiceName]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceName, []string{"CPU Min", "CPU Max", "CPU AVG", "CPU Total"})
		spark.SparkLines[metric.ServiceName] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "CPU Min" {
			sparkLines.Lines[i].Data = cpuMin
		} else if sparkLines.Lines[i].Title == "CPU Max" {
			sparkLines.Lines[i].Data = cpuMax
		} else if sparkLines.Lines[i].Title == "CPU AVG" {
			sparkLines.Lines[i].Data = cpuAvg
		} else if sparkLines.Lines[i].Title == "CPU Total" {
			sparkLines.Lines[i].Data = cpuTotal
		}
	}
	ui.Render(ui.Body)
}

func (spark *SparkTransformer) TransformSingleMemory(metric *models.Metrics) {
	var memMin []int
	var memMax []int
	var memAvg []int
	for _, data := range *metric.Data.MemoryUsage {
		memMin = append(memMin, int(data.Min/1024.0))
		memMax = append(memMax, int(data.Max/1024.0))
		memAvg = append(memAvg, int(data.AVG/1024.0))
	}
	var sparkLines = spark.SparkLines[metric.ServiceName]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceName, []string{"Mem Min", "Mem Max", "Mem AVG"})
		spark.SparkLines[metric.ServiceName] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "Mem Min" {
			sparkLines.Lines[i].Data = memMin
		} else if sparkLines.Lines[i].Title == "Mem Max" {
			sparkLines.Lines[i].Data = memMax
		} else if sparkLines.Lines[i].Title == "Mem AVG" {
			sparkLines.Lines[i].Data = memAvg
		}
	}
	ui.Render(ui.Body)
}

func (spark *SparkTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	var netinKB []int
	var netinPackets []int
	for _, data := range *metric.Data.NetworkUsage {
		netinKB = append(netinKB, int(data.RXKB))
		netinPackets = append(netinPackets, int(data.RXPackets))
	}
	var sparkLines = spark.SparkLines[metric.ServiceName]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceName, []string{"Received KB", "Received Packets"})
		spark.SparkLines[metric.ServiceName] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "Received KB" {
			sparkLines.Lines[i].Data = netinKB
		} else if sparkLines.Lines[i].Title == "Received Packets" {
			sparkLines.Lines[i].Data = netinPackets
		}
	}
	ui.Render(ui.Body)
}

func (spark *SparkTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	var netoutKB []int
	var netoutPackets []int
	for _, data := range *metric.Data.NetworkUsage {
		netoutKB = append(netoutKB, int(data.TXKB))
		netoutPackets = append(netoutPackets, int(data.TXPackets))
	}
	var sparkLines = spark.SparkLines[metric.ServiceName]
	if sparkLines == nil {
		sparkLines = addSparkLine(metric.ServiceName, []string{"Transmitted KB", "Transmitted Packets"})
		spark.SparkLines[metric.ServiceName] = sparkLines
	}
	for i := range sparkLines.Lines {
		if sparkLines.Lines[i].Title == "Transmitted KB" {
			sparkLines.Lines[i].Data = netoutKB
		} else if sparkLines.Lines[i].Title == "Transmitted Packets" {
			sparkLines.Lines[i].Data = netoutPackets
		}
	}
	ui.Render(ui.Body)
}

func addSparkLine(serviceName string, titles []string) *ui.Sparklines {
	var sparkLines []ui.Sparkline
	for _, title := range titles {
		sparkLine := ui.NewSparkline()
		sparkLine.Height = 1
		sparkLine.Data = []int{}
		sparkLine.Title = title
		sparkLines = append(sparkLines, sparkLine)
	}
	sp := ui.NewSparklines(sparkLines...)
	sp.Height = 14
	sp.BorderLabel = serviceName

	ui.Body.AddRows(
		ui.NewRow(ui.NewCol(12, 0, sp)),
	)

	ui.Body.Align()
	ui.Render(sp)
	ui.Render(ui.Body)

	return sp
}

func sparkLinesEventLoop() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Loop() // blocking call
}
