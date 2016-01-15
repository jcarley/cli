package metrics

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"strconv"

	"github.com/catalyzeio/cli/models"
)

// CSVTransformer is a concrete implementation of Transformer transforming
// data into CSV.
type CSVTransformer struct {
	HeadersWritten bool
	GroupMode      bool // I cant figure out a good way to access the Transformer.GroupMode from a CSVTransformer instance, so I copied it
	Buffer         *bytes.Buffer
	Writer         *csv.Writer
}

func (m *SMetrics) CSV() error {
	return nil
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
