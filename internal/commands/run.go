package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/rdaniels6813/cli-manager/internal/nodeman"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command using the installed CLI",
	Long: `Find an installed CLI application, and run a command by passing
 any subsequent arguments to that CLI instead`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		remainingArgs := args[1:]
		manager := nodeman.NewManager(afero.NewOsFs())
		commandPath, err := manager.GetCommandPath(command)
		if err != nil {
			log.Fatal(err)
		}
		runCommand := exec.Command(commandPath, remainingArgs...)
		manager.ConfigureNodeOnCommand(command, runCommand)
		runCommand.Stdout = os.Stdout
		runCommand.Stdin = os.Stdin
		runCommand.Stderr = os.Stderr
		err = runCommand.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.DisableFlagParsing = true
}
