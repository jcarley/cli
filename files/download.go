package files

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/catalyzeio/cli/helpers"
	"github.com/catalyzeio/cli/models"
)

// CmdDownload downloads a service file by its name (taken from listing
// service files) to the local machine matching the file's assigned permissions.
// If those permissions cannot be applied, the default 0644 permissions are
// applied. If not output file is specified, the file and permissions are
// printed to stdout.
func CmdDownload(ifiles IFiles, is services.IService) error {
	service, err := is.RetrieveByLabel()
	if err != nil {
		return err
	}
	//service := helpers.RetrieveServiceByLabel(serviceName, settings)
	if service == nil {
		return fmt.Errorf("Could not find a service with the name \"%s\"\n", serviceName)
	}
	if outputPath != "" && !force {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify `--force` to overwrite\n", outputPath)
		}
	}
	var file *models.ServiceFile
	files, err := ifiles.List()
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

func (f *SFiles) Retrieve() (*models.ServiceFile, error) {
	return nil, nil
}
