package main

import (
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

	intro, err := os.OpenFile("intro.md", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer intro.Close()
	output, err := os.OpenFile("cli-docs.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	_, err = io.Copy(output, intro)
	if err != nil {
		panic(err)
	}

	app.Cmd.DoInit()
	if _, err = output.WriteString("\n# Overview\n\n"); err != nil {
		panic(err)
	}
	output.Write([]byte("```\n"))
	app.Cmd.PrintLongHelpTo(false, output)
	output.Write([]byte("\n```\n\n"))
	for _, subCmd := range app.Cmd.Commands {
		if err := generateCommandDocs(subCmd, "#", output); err != nil {
			panic(err)
		}
	}
}

func generateCommandDocs(cmd *cli.Cmd, prefix string, writer io.Writer) error {
	cmd.DoInit()
	title := cmd.Name
	if len(title) > 1 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}
	writer.Write([]byte(prefix + " " + title + "\n\n"))
	if len(cmd.Commands) == 0 || cmd.LongDesc == "" { // shortcut to print out the Overview section
		writer.Write([]byte("```\n"))
		cmd.PrintLongHelpTo(false, writer)
		writer.Write([]byte("\n```\n\n"))
	}
	writer.Write([]byte(cmd.LongDesc))
	writer.Write([]byte("\n\n"))
	for _, subCmd := range cmd.Commands {
		if err := generateCommandDocs(subCmd, "#"+prefix+" "+title, writer); err != nil {
			return err
		}
	}
	return nil
}
