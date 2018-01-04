package test

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/daticahealth/cli/lib/prompts"
	"github.com/daticahealth/cli/models"
	notaryClient "github.com/docker/notary/client"
	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/data"
	"github.com/olekukonko/tablewriter"
)

// FakeImages is a mock images struct
type FakeImages struct {
	Settings *models.Settings
}

// Testing constants
const (
	Image = "hello-world"
	Tag   = "latest"
)

var testDigest = &models.ContentDigest{
	HashType: "sha256",
	Hash:     "8072a54ebb3bc136150e2f2860f00a7bf45f13eeb917cca2430fcd0054c8e51b",
	Size:     524,
}

var registries = map[string]string{}

var notaryServers = map[string]string{}

var localImages = []string{}

var remoteImages = []string{}

// AddRegistry adds a registry to the map
func AddRegistry(pod, registry string) {
	registries[pod] = registry
}

// AddNotary adds a notary server to the map
func AddNotary(pod, notary string) {
	notaryServers[pod] = notary
}

// SetLocalImages sets the mock list of local images
func SetLocalImages(images []string) {
	localImages = images
}

// SetRemoteImages sets the mock list of remote images
func SetRemoteImages(images []string) {
	remoteImages = images
}

// Errors for image handling
const (
	InvalidImageName             = "Invalid image name"
	IncorrectNamespace           = "Incorrect namespace for your environment"
	IncorrectRegistryOrNamespace = "Incorrect registry or namespace for your environment"
	MissingTrustData             = "does not have trust data for"
	ImageDoesNotExist            = "No such image"
)

// Constants for image handling
const (
	defaultTag = "latest"
	trustPath  = ".docker/trust"
)

// Push parses a given name into registry/namespace/image:tag and attempts to push it to the remote registry
func (d *FakeImages) Push(name string, user *models.User, env *models.Environment, ip prompts.IPrompts) (*models.Image, error) {
	repositoryName, tag, err := d.GetGloballyUniqueNamespace(name, env)
	if err != nil {
		return nil, err
	}

	if tag == "" {
		tag = defaultTag
	}
	fullImageName := strings.Join([]string{repositoryName, tag}, ":")

	if !imageExists(name, localImages) {
		imageError := fmt.Errorf(ImageDoesNotExist)
		if fullImageName != name && imageExists(fullImageName, localImages) {
			if yesNo := ip.YesNo(imageError.Error(), fmt.Sprintf("Would you like to push %s instead? (y/n) ", fullImageName)); yesNo != nil {
				return nil, imageError
			}
		} else {
			return nil, imageError
		}
	}

	return &models.Image{
		Name:   repositoryName,
		Tag:    tag,
		Digest: testDigest,
	}, nil
}

// Pull parses a name into registry/namespace/image:tag and attempts to retrieve it from the remote registry
func (d *FakeImages) Pull(name string, user *models.User, env *models.Environment) (*models.Image, error) {
	repositoryName, tag, err := d.GetGloballyUniqueNamespace(name, env)
	if err != nil {
		return nil, err
	}
	logrus.Printf("Pulling from repository %s\n", repositoryName)
	if tag == "" {
		tag = defaultTag
		logrus.Printf("Using default tag: %s\n", tag)
	}
	fullImageName := fmt.Sprintf("%s:%s", repositoryName, tag)
	if !imageExists(fullImageName, remoteImages) {
		return nil, fmt.Errorf(ImageDoesNotExist)
	}

	return &models.Image{
		Name:   repositoryName,
		Tag:    tag,
		Digest: testDigest,
	}, nil
}

// InitNotaryRepo intializes a notary repository
func (d *FakeImages) InitNotaryRepo(repo notaryClient.Repository, rootKeyPath string) error {
	rootTrustDir := fmt.Sprintf("%s/%s", userHomeDir(), trustPath)
	if err := os.MkdirAll(rootTrustDir, 0700); err != nil {
		return err
	}

	if err := repo.Initialize(nil); err != nil {
		return err
	}
	return nil
}

// AddTargetHash adds the given content hash to a notary repo and sends a signing request to the server
func (d *FakeImages) AddTargetHash(repo notaryClient.Repository, digest *models.ContentDigest, tag string, publish bool) error {
	targetHash := data.Hashes{}
	sha256, err := hex.DecodeString(digest.Hash)
	if err != nil {
		return err
	}
	targetHash[digest.HashType] = sha256

	// var targetCustom *canonicalJson.RawMessage
	target := &notaryClient.Target{Name: tag, Hashes: targetHash, Length: digest.Size}
	if err = repo.AddTarget(target, data.CanonicalTargetsRole); err != nil {
		return err
	}
	if publish {
		return d.Publish(repo)
	}
	return nil
}

// ListTargets intializes a notary repository
func (d *FakeImages) ListTargets(repo notaryClient.Repository, roles ...string) ([]*notaryClient.TargetWithRole, error) {
	target, _ := d.LookupTarget(repo, Tag)
	targets := []*notaryClient.TargetWithRole{target}
	return targets, nil
}

// LookupTarget searches for a specific target in a repository by tag name
func (d *FakeImages) LookupTarget(repo notaryClient.Repository, tag string) (*notaryClient.TargetWithRole, error) {
	hash, err := hex.DecodeString(testDigest.Hash)
	if err != nil {
		return nil, err
	}
	target := &notaryClient.TargetWithRole{
		Target: notaryClient.Target{
			Name:   tag,
			Hashes: data.Hashes{testDigest.HashType: hash},
			Length: testDigest.Size,
		},
		Role: "targets",
	}
	return target, nil
}

// DeleteTargets deletes the signed targets for a list of tags
func (d *FakeImages) DeleteTargets(repo notaryClient.Repository, tags []string, publish bool) error {
	return nil
}

// PrintChangelist prints out the users unpublished changes in a formatted table
func (d *FakeImages) PrintChangelist(changes []changelist.Change) {
	data := [][]string{[]string{"#", "Action", "Scope", "Type", "Target"}, []string{"-", "------", "-----", "----", "------"}}
	data = append(data)
	for i, c := range changes {
		data = append(data, []string{fmt.Sprintf("%d\n", i), c.Action(), c.Scope().String(), c.Type(), c.Path()})
	}

	logrus.Println()
	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetHeaderLine(false)
	table.SetAlignment(1)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()
	logrus.Println()
}

// CheckChangelist prompts the user if they have unpublished changes to let them clear
//	undesired changes before publishing, or verify that the changes should be published
func (d *FakeImages) CheckChangelist(repo notaryClient.Repository, ip prompts.IPrompts) error {
	if changelist, err := repo.GetChangelist(); err == nil {
		changes := changelist.List()
		if len(changes) > 0 {
			logrus.Println("The following unpublished changes were found in your local trust repository:")
			d.PrintChangelist(changes)
			publishWarning := "These changes will be published along with your current request."
			if err = ip.YesNo(publishWarning, "Would you like to proceed? (y/n) "); err != nil {
				logrus.Println("Use the `datica images targets reset <image>` command to clear undesired changes")
				return err
			}
		}
	}
	return nil
}

// GetNotaryRepository returns a pointer to the notary repository for an image
func (d *FakeImages) GetNotaryRepository(pod, imageName string, user *models.User) notaryClient.Repository {
	notaryServer := getServer(pod, notaryServers)
	rootTrustDir := fmt.Sprintf("%s/%s", userHomeDir(), trustPath)
	repo, err := notaryClient.NewFileCachedRepository(
		rootTrustDir,
		data.GUN(imageName),
		notaryServer,
		nil,
		nil,
		trustpinning.TrustPinConfig{},
	)
	if err != nil {
		logrus.Fatalln(err)
	}
	return repo
}

// GetGloballyUniqueNamespace returns the fully formatted name for an image <registry>/<namespace>/<image> and a tag if present
func (d *FakeImages) GetGloballyUniqueNamespace(name string, env *models.Environment) (string, string, error) {
	var repositoryName string
	var image string
	var tag string

	imageParts := strings.Split(name, ":")
	switch len(imageParts) {
	case 1:
		image = imageParts[0]
	case 2:
		image = imageParts[0]
		tag = imageParts[1]
	case 3:
		image = strings.Join([]string{imageParts[0], imageParts[1]}, ":")
		tag = imageParts[2]
	default:
		return "", "", fmt.Errorf(InvalidImageName)
	}

	repoParts := strings.Split(image, "/")
	registry := getServer(env.Pod, registries)

	switch len(repoParts) {
	case 1:
		repositoryName = fmt.Sprintf("%s/%s/%s", registry, env.Namespace, repoParts[0])
	case 2:
		if repoParts[0] != env.Namespace {
			if repoParts[0] != registry {
				return "", "", fmt.Errorf(IncorrectNamespace)
			}
			repositoryName = image
		} else {
			repositoryName = fmt.Sprintf("%s/%s", registry, image)
		}
	case 3:
		if repoParts[0] != registry || repoParts[1] != env.Namespace {
			return "", "", fmt.Errorf(IncorrectRegistryOrNamespace)
		}
		repositoryName = image
	default:
		return "", "", fmt.Errorf(InvalidImageName)
	}
	return repositoryName, tag, nil
}

// Publish pretends to publish changes. But actually just clear the changes
func (d *FakeImages) Publish(repo notaryClient.Repository) error {
	changelist, err := repo.GetChangelist()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return changelist.Clear("")
}

// DeleteLocalRepo deletes the local trust repo created in testing
func DeleteLocalRepo(namespace, image, registry, notaryServer string) error {
	repositoryName := strings.Join([]string{registry, namespace, image}, "/")
	rootTrustDir := fmt.Sprintf("%s/%s", userHomeDir(), trustPath)
	if err := notaryClient.DeleteTrustData(
		rootTrustDir,
		data.GUN(repositoryName),
		notaryServer,
		nil,
		false,
	); err != nil {
		return err
	}

	//Delete parent directories only if empty (just in case)
	repoDir := strings.Join([]string{rootTrustDir, "tuf", registry}, "/")
	namespaceDir := strings.Join([]string{repoDir, namespace}, "/")
	if err := os.Remove(namespaceDir); err != nil {
		fmt.Printf("Unable to remove directory %s:\n		%s\n", namespaceDir, err.Error())
		return nil
	}
	if err := os.Remove(repoDir); err != nil {
		fmt.Printf("Unable to remove directory %s:\n		%s\n", repoDir, err.Error())
	}
	return nil
}

// Fake image verification
func imageExists(image string, imageList []string) bool {
	for _, img := range imageList {
		if image == img {
			return true
		}
	}
	return false
}

func getServer(pod string, serverMap map[string]string) string {
	if server, ok := serverMap[pod]; ok {
		return server
	}
	return serverMap["default"]
}

func userHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

// ListImages stub to make golinter happy
func (d *FakeImages) ListImages() (*[]string, error) {
	return nil, nil
}

// ListTags stub to make golinter happy
func (d *FakeImages) ListTags(imageName string) (*[]string, error) {
	return nil, nil
}

// DeleteTag stub to make golinter happy
func (d *FakeImages) DeleteTag(imageName, tagName string) error {
	return nil
}
