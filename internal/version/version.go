package version

import (
	"fmt"
	"runtime/debug"
)

var version = ""

func GetModuleVersion() string {
	if version != "" {
		return version
	}
	if bi, exists := debug.ReadBuildInfo(); exists {
		return bi.Main.Version
	} else {
		return fmt.Sprintf("No version information found. Make sure to use " +
			"GO111MODULE=on when running 'go get' in order to use specific " +
			"version of the binary.")
	}

}
