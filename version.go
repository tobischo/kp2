package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "WORK IN PROGRESS"

func versionCmd(_ *cobra.Command, _ []string) {
	fmt.Println(app, "version", version)
}
