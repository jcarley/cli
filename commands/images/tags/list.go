package tags

import (
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/commands/environments"
	"github.com/daticahealth/cli/lib/images"
)

func cmdTagList(ii images.IImages, ie environments.IEnvironments, envID, image string) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	namespacedImage, _, err := ii.GetGloballyUniqueNamespace(image, env, false)
	if err != nil {
		return err
	}
	tags, err := ii.ListTags(namespacedImage)
	if err != nil {
		return err
	}
	logrus.Printf("Tags for image \"%s\"", image)
	logrus.Println("")
	sort.Strings(*tags)
	for _, tag := range *tags {
		logrus.Println(tag)
	}
	return nil
}
