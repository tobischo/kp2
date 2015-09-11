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

func loadDatabaseCmd(cmd *cobra.Command, args []string) error {
	var password string
	if usePassword {
		fmt.Println("Enter password:")
		pw, err := terminal.ReadPassword(0)
		if err != nil {
			return fmt.Errorf("Failed to read password: '%s'", err)
		}

		password = string(pw)
	}

	var credentials *gokeepasslib.DBCredentials
	var err error

	switch {
	case usePassword && keyFile != "":
		credentials, err = gokeepasslib.NewPasswordAndKeyCredentials(
			password, keyFile,
		)
	case usePassword:
		credentials = gokeepasslib.NewPasswordCredentials(
			password,
		)
	case keyFile != "":
		credentials, err = gokeepasslib.NewKeyCredentials(keyFile)
	default:
		return fmt.Errorf("Key file or password has to be provided")
	}

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

	return nil
}
