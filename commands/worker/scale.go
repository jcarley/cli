package worker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
)

func CmdScale(svcName, target, scaleString string, iw IWorker, is services.IServices, ip prompts.IPrompts, ij jobs.IJobs) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	scaleFunc, changeInScale, err := iw.ParseScale(scaleString)
	if err != nil {
		return err
	}
	workers, err := iw.Retrieve(service.ID)
	if err != nil {
		return err
	}
	scale := scaleFunc(workers.Workers[target], changeInScale)
	if existingScale, ok := workers.Workers[target]; !ok || scale > existingScale {
		logrus.Printf("Deploying %d new workers with target %s for service %s", scale-existingScale, target, svcName)
		workers.Workers[target] = scale
		err = iw.Update(service.ID, workers)
		if err != nil {
			return err
		}
		err = ij.DeployTarget(target, service.ID)
		if err != nil {
			return err
		}
		logrus.Printf("Successfully deployed %d new workers with target %s for service %s and set the scale to %d", scale-existingScale, target, svcName, scale)
	} else if scale < existingScale {
		err = ip.YesNo(fmt.Sprintf("Scaling down the %s target from %d to %d for service %s will automatically stop %d jobs, would you like to proceed? (y/n) ", target, existingScale, scale, svcName, existingScale-scale))
		if err != nil {
			return err
		}
		jobs, err := ij.RetrieveByTarget(service.ID, target, 1, 1000)
		if err != nil {
			return err
		}
		deleteLimit := existingScale - scale
		deleted := 0

		for _, j := range *jobs {
			err = ij.Delete(j.ID, service.ID)
			if err != nil {
				return err
			}
			deleted++
			if deleted == deleteLimit {
				break
			}
		}
		workers.Workers[target] = scale
		err = iw.Update(service.ID, workers)
		if err != nil {
			return err
		}
		logrus.Printf("Successfully removed %d existing workers with target %s for service %s and set the scale to %d", existingScale-scale, target, svcName, scale)
	} else {
		logrus.Printf("Worker target %s for service %s is already at a scale of %d", target, svcName, scale)
	}
	return nil
}

func (w *SWorker) Update(svcID string, workers *models.Workers) error {
	b, err := json.Marshal(workers)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(w.Settings.SessionToken, w.Settings.Version, w.Settings.Pod, w.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/workers", w.Settings.PaasHost, w.Settings.PaasHostVersion, w.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}

func (w *SWorker) ParseScale(scaleString string) (func(scale, change int) int, int, error) {
	scale, err := strconv.Atoi(scaleString)
	if err != nil {
		return nil, 0, err
	}

	if strings.HasPrefix(scaleString, "+") || strings.HasPrefix(scaleString, "-") {
		return changeScale, scale, nil
	}
	return constantScale, scale, nil
}

func changeScale(scale, change int) int {
	return scale + change
}

func constantScale(scale, newScale int) int {
	return newScale
}
