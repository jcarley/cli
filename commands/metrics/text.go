package metrics

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// TextTransformer is a concrete implementation of Transformer transforming
// data into plain text.
type TextTransformer struct{}

// TransformGroup transforms an entire environment's metrics data into text
// format.
func (text *TextTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		logrus.Printf("%s:", metric.ServiceName)
		text.TransformSingleCPU(&metric)
	}
}

func (text *TextTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		logrus.Printf("%s:", metric.ServiceName)
		text.TransformSingleMemory(&metric)
	}
}

func (text *TextTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		logrus.Printf("%s:", metric.ServiceName)
		text.TransformSingleNetworkIn(&metric)
	}
}

func (text *TextTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		logrus.Printf("%s:", metric.ServiceName)
		text.TransformSingleNetworkOut(&metric)
	}
}

// TransformSingle transforms a single service's metrics data into text
// format.
func (text *TextTransformer) TransformSingleCPU(metric *models.Metrics) {
	prefix := "    "
	for _, data := range *metric.Data.CPULoad {
		ts := time.Unix(int64(data.TS), 0)
		logrus.Printf("%s%s | CPU Min: %6.2fs | CPU Max: %6.2fs | CPU AVG: %6.2fs | CPU Total: %6.2fs",
			prefix,
			fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
			data.Min/1000.0,
			data.Max/1000.0,
			data.AVG/1000.0,
			data.Total/1000.0)
	}
}

func (text *TextTransformer) TransformSingleMemory(metric *models.Metrics) {
	prefix := "    "
	for _, data := range *metric.Data.MemoryUsage {
		ts := time.Unix(int64(data.TS), 0)
		logrus.Printf("%s%s | Memory Min: %.2f KB | Memory Max: %.2f KB | Memory AVG: %.2f KB",
			prefix,
			fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
			data.Min/1024.0,
			data.Max/1024.0,
			data.AVG/1024.0)
	}
}

func (text *TextTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	prefix := "    "
	for _, data := range *metric.Data.NetworkUsage {
		ts := time.Unix(int64(data.TS), 0)
		logrus.Printf("%s%s | Received Bytes: %.2f KB | Received Packets: %.2f",
			prefix,
			fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
			data.RXKB,
			data.RXPackets)
	}
}

func (text *TextTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	prefix := "    "
	for _, data := range *metric.Data.NetworkUsage {
		ts := time.Unix(int64(data.TS), 0)
		logrus.Printf("%s%s | Transmitted Bytes: %.2f KB | Transmitted Packets: %.2f",
			prefix,
			fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second()),
			data.TXKB,
			data.TXPackets)
	}
}
