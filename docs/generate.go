package main

import (
	"io"
	"os"
	"strings"

	"github.com/catalyzeio/cli/catalyze"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

const binaryName = "catalyze"

func main() {
	app := cli.App("catalyze", "")
	settings := &models.Settings{}
	catalyze.InitLogrus()
	catalyze.InitGlobalOpts(app, settings)
	catalyze.InitCLI(app, settings)

	output, err := os.OpenFile("cli-docs.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	app.Cmd.DoInit()
	if _, err = output.WriteString("# Overview\n\n"); err != nil {
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
