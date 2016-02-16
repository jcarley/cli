package metrics

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

// JSONTransformer is a concrete implementation of Transformer transforming
// data into JSON.
type JSONTransformer struct{}

type cpu struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	AVG         float64 `json:"avg"`
	Total       float64 `json:"total"`
}

type mem struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	AVG         float64 `json:"avg"`
}

type netin struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	RXKB        float64 `json:"rx_kb"`
	RXPackets   float64 `json:"rx_packets"`
}

type netout struct {
	ServiceName string  `json:"service_name,omitempty"`
	TS          int     `json:"ts"`
	TXKB        float64 `json:"tx_kb"`
	TXPackets   float64 `json:"tx_packets"`
}

// TransformGroup transforms an entire environment's metrics data into json
// format.
func (j *JSONTransformer) TransformGroupCPU(metrics *[]models.Metrics) {
	var data []cpu
	for _, m := range *metrics {
		for _, d := range *m.Data.CPULoad {
			data = append(data, cpu{m.ServiceLabel, d.TS, d.Min / 1000.0, d.Max / 1000.0, d.AVG / 1000.0, d.Total / 1000.0})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

func (j *JSONTransformer) TransformGroupMemory(metrics *[]models.Metrics) {
	var data []mem
	for _, m := range *metrics {
		for _, d := range *m.Data.MemoryUsage {
			data = append(data, mem{m.ServiceLabel, d.TS, d.Min / 1024.0, d.Max / 1024.0, d.AVG / 1024.0})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

func (j *JSONTransformer) TransformGroupNetworkIn(metrics *[]models.Metrics) {
	var data []netin
	for _, m := range *metrics {
		for _, d := range *m.Data.NetworkUsage {
			data = append(data, netin{m.ServiceLabel, d.TS, d.RXKB, d.RXPackets})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

func (j *JSONTransformer) TransformGroupNetworkOut(metrics *[]models.Metrics) {
	var data []netout
	for _, m := range *metrics {
		for _, d := range *m.Data.NetworkUsage {
			data = append(data, netout{m.ServiceLabel, d.TS, d.TXKB, d.TXPackets})
		}
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

// TransformSingle transforms a single service's metrics data into json
// format.
func (j *JSONTransformer) TransformSingleCPU(metric *models.Metrics) {
	var data []cpu
	for _, d := range *metric.Data.CPULoad {
		data = append(data, cpu{TS: d.TS, Min: d.Min / 1000.0, Max: d.Max / 1000.0, AVG: d.AVG / 1000.0, Total: d.Total / 1000.0})
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

func (j *JSONTransformer) TransformSingleMemory(metric *models.Metrics) {
	var data []mem
	for _, d := range *metric.Data.MemoryUsage {
		data = append(data, mem{TS: d.TS, Min: d.Min / 1024.0, Max: d.Max / 1024.0, AVG: d.AVG / 1024.0})
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

func (j *JSONTransformer) TransformSingleNetworkIn(metric *models.Metrics) {
	var data []netin
	for _, d := range *metric.Data.NetworkUsage {
		data = append(data, netin{TS: d.TS, RXKB: d.RXKB, RXPackets: d.RXPackets})
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}

func (j *JSONTransformer) TransformSingleNetworkOut(metric *models.Metrics) {
	var data []netout
	for _, d := range *metric.Data.NetworkUsage {
		data = append(data, netout{TS: d.TS, TXKB: d.TXKB, TXPackets: d.TXPackets})
	}
	b, _ := json.MarshalIndent(data, "", "    ")
	logrus.Println(string(b))
}
