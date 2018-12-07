package main

import (
	"fmt"

	"github.com/tobischo/gokeepasslib/v2"
	"golang.org/x/crypto/ssh/terminal"
)

func readString(text string) (string, error) {
	fmt.Print(text)
	var input string

	_, err := fmt.Scanln(&input)

	return input, err
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

func readPasswordWithConfirmation() (string, error) {
	var (
		password         string
		passwordRepeated string
		err              error
	)

	password, err = readPassword("Enter password: ")
	if err != nil {
		return "", err
	}

	passwordRepeated, err = readPassword("Repeat password: ")
	if err != nil {
		return "", err
	}

	if password != passwordRepeated {
		return "", fmt.Errorf("Password and repetition do not match")
	}

	return password, nil
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
