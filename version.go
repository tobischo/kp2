package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "WORK IN PROGRESS"

func versionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(app, "version", version)
}
