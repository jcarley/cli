package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/daticahealth/cli/datica"
	"github.com/daticahealth/cli/models"
	"github.com/jault3/mow.cli"
)

const binaryName = "datica"

func main() {
	app := cli.App(binaryName, "")
	settings := &models.Settings{}
	datica.InitLogrus()
	datica.InitGlobalOpts(app, settings)
	datica.InitCLI(app, settings)

	intro, err := os.OpenFile("intro.html", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer intro.Close()
	output, err := os.OpenFile("cli-docs.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	_, err = io.Copy(output, intro)
	if err != nil {
		panic(err)
	}

	app.Cmd.DoInit()
	if _, err = output.WriteString("\n<h1>Overview</h1>\n\n"); err != nil {
		panic(err)
	}
	output.Write([]byte("<pre>\n"))
	app.Cmd.PrintLongHelpTo(false, output)
	output.Write([]byte("\n</pre>\n\n"))
	for _, subCmd := range app.Cmd.Commands {
		if err := generateCommandDocs(subCmd, 2, output); err != nil {
			panic(err)
		}
	}
}

func generateCommandDocs(cmd *cli.Cmd, tagLevel int, writer io.Writer) error {
	cmd.DoInit()
	title := cmd.Name
	if len(title) > 1 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}
	writer.Write([]byte(fmt.Sprintf("<h%[1]d>%[2]s</h%[1]d>\n\n", tagLevel, title)))
	if len(cmd.Commands) == 0 || cmd.LongDesc == "" { // shortcut to print out the Overview section
		writer.Write([]byte("<pre>\n"))
		cmd.PrintLongHelpTo(false, writer)
		writer.Write([]byte("\n</pre>\n\n"))
	}
	writer.Write([]byte(cmd.LongDesc))
	writer.Write([]byte("\n\n"))
	for _, subCmd := range cmd.Commands {
		if err := generateCommandDocs(subCmd, tagLevel+1, writer); err != nil {
			return err
		}
	}
	return nil
}
