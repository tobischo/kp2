package main

import (
	"os"

	"github.com/spf13/cobra"
)

func paramSetupCmd(cmd *cobra.Command, args []string) {
	if file == "" {
		file = os.Getenv("KP2FILE")
	}

	if keyFile == "" {
		usePassword = true
	}
}
