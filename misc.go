package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/jakewarren/trustar-golang"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

// command to list available enclaves
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists enclaves",
	Run: func(cmd *cobra.Command, args []string) {
		// get Enclaves
		enclaves, err := c.GetEnclaves()
		if err != nil {
			fmt.Println(err)
		}

		// TODO: add option to output as json?

		prettyPrint(enclaves)
	},
}

// command to print token
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print access token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(c.Token.Token)
	},
}

func prettyPrint(enclaves []trustar.Enclave) {
	headerFmt := color.New(color.FgCyan, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("id", "name", "type", "read", "update", "create")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, e := range enclaves {
		tbl.AddRow(e.ID, e.Name, e.Type, e.Read, e.Update, e.Create)
	}
	tbl.Print()
}

// TODO: document how to use this in the README
// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "autocomplete",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(trustar autocomplete)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(trustar autocomplete)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

// openBrowser opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
