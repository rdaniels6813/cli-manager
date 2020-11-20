package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

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
			fmt.Println(err)
			os.Exit(1)
		}
		c := exec.Command(commandPath, remainingArgs...)
		manager.ConfigureNodeOnCommand(command, c)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		if err := c.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		done := false
		notif := make(chan os.Signal, 1)
		signal.Notify(notif, os.Interrupt)
		go func() {
			for range notif {
				if done {
					os.Exit(0)
				}
			}
		}()
		err = c.Wait()
		done = true
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.DisableFlagParsing = true
}
