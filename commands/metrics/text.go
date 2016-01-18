package metrics

import (
	"fmt"
	"math"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// TextTransformer is a concrete implementation of Transformer transforming
// data into plain text.
type TextTransformer struct{}

func (m *SMetrics) Text() error {
	return nil
}

// TransformGroup transforms an entire environment's metrics data into text
// format.
func (text *TextTransformer) TransformGroup(metrics *[]models.Metrics) {
	for _, metric := range *metrics {
		logrus.Printf("%s:", metric.ServiceName)
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
			logrus.Printf("%s%s | %8s (%s) | CPU: %6.2fs (%5.2f%%) | Net: RX: %.2f KB TX: %.2f KB | Mem: %.2f KB | Disk: %.2f KB read / %.2f KB write",
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
