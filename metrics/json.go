package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/catalyzeio/cli/models"
)

// JSONTransformer is a concrete implementation of Transformer transforming
// data into JSON.
type JSONTransformer struct{}

func (m *SMetrics) JSON() error {
	return nil
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
