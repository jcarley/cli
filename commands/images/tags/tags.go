package tags

import (
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/lib/images"
)

func cmdTagList(ii images.IImages, image string) error {
	tags, err := ii.ListTags(image)
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
