package tags

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/prompts"
)

func cmdTagDelete(ii images.IImages, ip prompts.IPrompts, ie environments.IEnvironments, envID, image string) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	namespacedImage, tag, err := ii.GetGloballyUniqueNamespace(image, env, false)
	if err != nil {
		return err
	} else if tag == "" {
		return fmt.Errorf("Must include tag in image name.")
	}

	tags, err := ii.ListTags(namespacedImage)
	if err != nil {
		return err
	}
	tagFound := false
	for _, t := range *tags {
		if tag == t {
			tagFound = true
			break
		}
	}
	if !tagFound {
		return errors.New("No tags found matching the given name for the given image on this environment.")
	} else {
		logrus.Println("Warning! Deleting a tag will also delete any other tags that point to identical images.")
		err = ip.YesNo("", "Are you sure you want to delete this tag? (y/n) ")
		if err != nil {
			return err
		}
		err = ii.DeleteTag(namespacedImage, tag)
		if err == nil {
			logrus.Println("Tag deleted successfully.")
		}
		return err
	}
}
