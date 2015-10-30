package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib"
)

func paramSetupCmd(cmd *cobra.Command, args []string) {
	if filePath == "" {
		filePath = os.Getenv("KP2FILE")
	}

	if keyFile == "" {
		usePassword = true
	}
}

func readPassword(text string) (string, error) {
	fmt.Print(text)
	pw, err := terminal.ReadPassword(0)
	if err != nil {
		return "", fmt.Errorf("Failed to read password: '%s'", err)
	}
	fmt.Println()
	return string(pw), nil
}

func pickCredentialMode(password string) (*gokeepasslib.DBCredentials, error) {
	switch {
	case usePassword && keyFile != "":
		return gokeepasslib.NewPasswordAndKeyCredentials(
			password, keyFile,
		)
	case usePassword:
		credentials := gokeepasslib.NewPasswordCredentials(
			password,
		)
		return credentials, nil
	case keyFile != "":
		return gokeepasslib.NewKeyCredentials(keyFile)
	default:
		return nil, fmt.Errorf("Key file or password has to be provided")
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
