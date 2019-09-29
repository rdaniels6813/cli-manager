package main

import cmd "github.com/rdaniels6813/cli-manager/pkg/commands"

var Version string

func main() {
	cmd.RootCmd.Version = Version
	cmd.Execute()
}
