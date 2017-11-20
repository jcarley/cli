package tags

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/lib/images"
	"github.com/daticahealth/cli/lib/prompts"
)

func cmdTagDelete(ii images.IImages, ip prompts.IPrompts, image, tagName string) error {
	tags, err := ii.ListTags(image)
	if err != nil {
		return err
	}
	tagFound := false
	for _, tag := range *tags {
		if tag == tagName {
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
		err = ii.DeleteTag(image, tagName)
		if err == nil {
			logrus.Println("Tag deleted successfully.")
		}
		return err
	}
}
