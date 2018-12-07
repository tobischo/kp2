package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib/v2"
)

func paramSetupCmd(cmd *cobra.Command, args []string) {
	if filePath == "" {
		filePath = os.Getenv("KP2FILE")
	}

	if keyFile == "" {
		usePassword = true
	}
}

func loadDatabaseCmd(cmd *cobra.Command, args []string) error {
	var (
		password string
		err      error
	)

	if usePassword {
		password, err = readPassword("Enter password: ")
		if err != nil {
			return err
		}
	}

	credentials, err := pickCredentialMode(password)
	if err != nil {
		return fmt.Errorf("Failed to setup credentials: '%s'", err)
	}

	db = new(gokeepasslib.Database)
	db.Credentials = credentials

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open Keepass2 file %s: '%s'", filePath, err)
	}

	err = gokeepasslib.NewDecoder(file).Decode(db)
	if err != nil {
		return fmt.Errorf("Failed to decode Keepass2 file: %s", err)
	}

	if err := db.UnlockProtectedEntries(); err != nil {
		return err
	}

	return nil
}
