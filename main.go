package main

import (
	"log"

	"lab.bittrd.com/bittrd/cli-manager/nodeman"

	"github.com/spf13/afero"
)

func main() {
	nodeManager := nodeman.NewManager(afero.NewOsFs())
	node := nodeManager.GetNode("10.16.0")
	err := node.Node("-v")
	log.Println(err)
	err = node.Npm("-v")
	log.Println(err)
}
