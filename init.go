package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib/v3"
)

func initCmd(_ *cobra.Command, _ []string) error {
	db = gokeepasslib.NewDatabase()

	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("File at '%s' already exists: %w", filePath, err)
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
		return fmt.Errorf("failed to setup credentials: '%w'", err)
	}

	db.Credentials = credentials
	changed = true

	return nil
}
