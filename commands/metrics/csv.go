package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// CSVTransformer is a concrete implementation of Transformer transforming
// data into CSV.
type CSVTransformer struct {
	HeadersWritten bool
	GroupMode      bool
	Buffer         *bytes.Buffer
	Writer         *csv.Writer
}

// WriteHeaders prints out the csv headers to the console.
func (csv *CSVTransformer) WriteHeadersCPU() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "cpu_min", "cpu_max", "cpu_avg", "cpu_total"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

func (csv *CSVTransformer) WriteHeadersMemory() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "memory_min", "memory_max", "memory_avg"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

func (csv *CSVTransformer) WriteHeadersNetworkIn() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "rx_bytes", "rx_packets"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

func (csv *CSVTransformer) WriteHeadersNetworkOut() {
	if !csv.HeadersWritten {
		headers := []string{"timestamp", "tx_bytes", "tx_packets"}
		if csv.GroupMode {
			headers = append([]string{"service_name"}, headers...)
		}
		csv.Writer.Write(headers)
		csv.HeadersWritten = true
	}
}

// TransformGroup transforms an entire environment's metrics data into csv
// format.
func (csv *CSVTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		csv.TransformSingleCPU(&metric)
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

func (csv *CSVTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		csv.TransformSingleMemory(&metric)
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

func (csv *CSVTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		csv.TransformSingleNetworkIn(&metric)
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

func (csv *CSVTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	csv.GroupMode = true
	for _, metric := range *metrics {
		csv.TransformSingleNetworkOut(&metric)
	}
	csv.Writer.Flush()
	logrus.Println(csv.Buffer.String())
}

// TransformSingle transforms a single service's metrics data into csv
// format.
func (csv *CSVTransformer) TransformSingleCPU(metric *models.Metrics) {
	csv.WriteHeadersCPU()
	for _, data := range *metric.Data.CPULoad {
		row := []string{
			fmt.Sprintf("%d", data.TS),
			fmt.Sprintf("%f", data.Min/1000.0),
			fmt.Sprintf("%f", data.Max/1000.0),
			fmt.Sprintf("%f", data.AVG/1000.0),
			fmt.Sprintf("%f", data.Total/1000.0),
		}
		if csv.GroupMode {
			row = append([]string{metric.ServiceName}, row...)
		}
		csv.Writer.Write(row)
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}

func (csv *CSVTransformer) TransformSingleMemory(metric *models.Metrics) {
	csv.WriteHeadersMemory()
	for _, data := range *metric.Data.MemoryUsage {
		row := []string{
			fmt.Sprintf("%d", data.TS),
			fmt.Sprintf("%f", data.Min/1024.0),
			fmt.Sprintf("%f", data.Max/1024.0),
			fmt.Sprintf("%f", data.AVG/1024.0),
		}
		if csv.GroupMode {
			row = append([]string{metric.ServiceName}, row...)
		}
		csv.Writer.Write(row)
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}

func (csv *CSVTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	csv.WriteHeadersNetworkIn()
	for _, data := range *metric.Data.NetworkUsage {
		row := []string{
			fmt.Sprintf("%d", data.TS),
			fmt.Sprintf("%f", data.RXKB),
			fmt.Sprintf("%f", data.RXPackets),
		}
		if csv.GroupMode {
			row = append([]string{metric.ServiceName}, row...)
		}
		csv.Writer.Write(row)
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}

func (csv *CSVTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	csv.WriteHeadersNetworkOut()
	for _, data := range *metric.Data.NetworkUsage {
		row := []string{
			fmt.Sprintf("%d", data.TS),
			fmt.Sprintf("%f", data.TXKB),
			fmt.Sprintf("%f", data.TXPackets),
		}
		if csv.GroupMode {
			row = append([]string{metric.ServiceName}, row...)
		}
		csv.Writer.Write(row)
	}
	if !csv.GroupMode {
		csv.Writer.Flush()
		logrus.Println(csv.Buffer.String())
	}
}
