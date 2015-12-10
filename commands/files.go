package commands

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
)

// ListServiceFiles lists all service files that are able to be downloaded
// by a member of the environment. Typically service files of interest
// will be on the service_proxy.
func ListServiceFiles(serviceName string, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceName, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the name \"%s\"\n", serviceName)
		os.Exit(1)
	}
	files := helpers.ListServiceFiles(service.ID, settings)
	if len(*files) == 0 {
		fmt.Println("No service files found")
		return
	}
	fmt.Println("NAME")
	for _, sf := range *files {
		fmt.Println(sf.Name)
	}
}

// DownloadServiceFile downloads a service file by its ID (take from listing
// service files) to the local machine matching the file's assigned permissions.
// If those permissions cannot be applied, the default 0644 permissions are
// applied. If not output file is specified, the file and permissions are
// printed to stdout.
func DownloadServiceFile(serviceName, fileName, outputPath string, force bool, settings *models.Settings) {
	helpers.SignIn(settings)
	service := helpers.RetrieveServiceByLabel(serviceName, settings)
	if service == nil {
		fmt.Printf("Could not find a service with the name \"%s\"\n", serviceName)
		os.Exit(1)
	}
	if outputPath != "" && !force {
		if _, err := os.Stat(outputPath); err == nil {
			fmt.Printf("File already exists at path '%s'. Specify `--force` to overwrite\n", outputPath)
			os.Exit(1)
		}
	}
	var file *models.ServiceFile
	files := helpers.ListServiceFiles(service.ID, settings)
	for _, f := range *files {
		if f.Name == fileName {
			file = helpers.RetrieveServiceFile(service.ID, f.ID, settings)
			break
		}
	}
	if file == nil {
		fmt.Printf("File with name %s does not exist. Try listing files again by running \"catalyze files list\"\n", fileName)
		os.Exit(1)
	}
	filePerms, err := strconv.ParseUint(file.Mode, 8, 32)
	if err != nil {
		filePerms = 0644
	}

	var wr io.Writer
	if outputPath != "" {
		if force {
			os.Remove(outputPath)
		}
		outFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, os.FileMode(filePerms))
		if err != nil {
			fmt.Printf("Warning! Could not apply %s file permissions. Attempting to apply defaults %s\n", fileModeToRWXString(filePerms), fileModeToRWXString(uint64(0644)))
			outFile, err = os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				fmt.Printf("Could not open %s for writing: %s\n", outputPath, err.Error())
				os.Exit(1)
			}
		}
		defer outFile.Close()
		wr = outFile
	} else {
		fmt.Printf("Mode: %s\n\nContent:\n", fileModeToRWXString(filePerms))
		wr = os.Stdout
	}
	wr.Write([]byte(file.Contents))
}

func fileModeToRWXString(perms uint64) string {
	permissionString := ""
	binaryString := strconv.FormatUint(perms, 2)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if string(binaryString[len(binaryString)-1-i*3-j]) == "1" {
				switch j {
				case 0:
					permissionString = "x" + permissionString
				case 1:
					permissionString = "w" + permissionString
				case 2:
					permissionString = "r" + permissionString
				}
			} else {
				permissionString = "-" + permissionString
			}
		}
	}
	permissionString = "-" + permissionString // we don't store folders
	return permissionString
}
