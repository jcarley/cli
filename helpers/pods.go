package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/catalyzeio/catalyze/httpclient"
	"github.com/catalyzeio/catalyze/models"
)

// RetrievePodMetadata retrieves information about a certain Pod API
func RetrievePodMetadata(podID string, settings *models.Settings) *models.PodMetadata {
	resp := httpclient.Get(fmt.Sprintf("%s/v1/pods/metadata", settings.PaasHost), true, settings)
	var pods []models.PodMetadata
	json.Unmarshal(resp, &pods)
	var pod models.PodMetadata
	for _, p := range pods {
		if p.ID == podID {
			pod = p
			break
		}
	}
	if pod.ID == "" {
		fmt.Println("Could not find the pod associated with your environment. Please contact Catalyze support (support@catalyze.io). Please include your environment ID - found via \"catalyze support-ids\"")
		os.Exit(1)
	}
	return &pod
}
