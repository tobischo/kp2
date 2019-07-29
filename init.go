package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib/v3"
)

func initCmd(cmd *cobra.Command, args []string) error {
	db = gokeepasslib.NewDatabase()

	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("File at '%s' already exists", filePath)
	}

	var (
		password string
		err      error
	)

	if usePassword {
		password, err = readPasswordWithConfirmation()
		if err != nil {
			return err
		}
	}

	credentials, err := pickCredentialMode(password)
	if err != nil {
		return fmt.Errorf("Failed to setup credentials: '%s'", err)
	}

	db.Credentials = credentials
	changed = true

	return nil
}
