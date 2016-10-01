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
	app := cli.App("", "")
	settings := &models.Settings{}
	catalyze.InitLogrus()
	catalyze.InitGlobalOpts(app, settings)
	catalyze.InitCLI(app, settings)

	output, err := os.OpenFile("cli-docs.md", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	if _, err = output.WriteString("# Overview"); err != nil {
		panic(err)
	}

	if err = generateCommandDocs(app.Cmd, "", output); err != nil {
		panic(err)
	}
}

func generateCommandDocs(cmd *cli.Cmd, prefix string, writer io.Writer) error {
	cmd.DoInit()
	/*c := reflect.ValueOf(*cmd)
	name := c.FieldByName("name").String()*/
	title := cmd.Name
	if len(title) > 1 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}
	writer.Write([]byte(prefix + " " + title + "\n\n"))
	/*args := append(parents, cmd.Name)
	args = append(args, "--help")
	output := runCommand(binaryName, args...)
	_, err := writer.Write(bytes.TrimSpace(output))
	if err != nil {
		return err
	}*/
	if len(cmd.Commands) == 0 || cmd.Name == "" { // shortcut to print out the Overview section
		writer.Write([]byte("```\n"))
		cmd.PrintLongHelpTo(false, writer)
		writer.Write([]byte("\n```\n\n"))
	}
	writer.Write([]byte(cmd.LongDesc))
	writer.Write([]byte("\n\n"))
	// cmdsFields := c.FieldByName("commands")
	// for i := 0; i < cmdsFields.Len(); i++ {
	for _, subCmd := range cmd.Commands {
		// // cmdElem := cmdsFields.Index(i).Elem()
		// // ptr := unsafe.Pointer(cmdElem.UnsafeAddr())
		// // subCmd := (*cli.Cmd)(ptr)
		// // args := append(parents, name)
		// err = generateCommandDocs(subCmd, "#"+prefix+" "+title, writer, args...)
		if err := generateCommandDocs(subCmd, "#"+prefix+" "+title, writer); err != nil {
			return err
		}
	}
	return nil
}

// func runCommand(cmdName string, args ...string) []byte {
// 	fmt.Printf("Running %s %v\n", cmdName, args)
// 	stdout := &bytes.Buffer{}
// 	stderr := &bytes.Buffer{}
// 	cmd := exec.Command(cmdName, args...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = stdout
// 	cmd.Env = os.Environ()
// 	cmd.Stderr = stderr
// 	cmd.Run()
// 	return stderr.Bytes()
// }
