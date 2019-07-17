// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/afero"
	"lab.bittrd.com/bittrd/cli-manager/nodeman"

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
