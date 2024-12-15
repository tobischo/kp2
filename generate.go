package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

const (
	lowerCaseBytes = "abcdefghijklmnopqrstuvwxyz"
	upperCaseBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitBytes     = "0123456789"
	specialBytes   = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

func generateCmd(_ *cobra.Command, _ []string) error {
	var passwordBytes string

	passwordLengthStr, err := readString("Desired password length? ")
	if err != nil {
		return err
	}

	passwordLength, err := strconv.Atoi(passwordLengthStr)
	if err != nil {
		return err
	}

	addLowerCase, err := readString("Lower case letters? [y/n] ")
	if err != nil {
		return err
	}

	if addLowerCase == `y` {
		passwordBytes += lowerCaseBytes
	}

	addUpperCase, err := readString("Upper case letter? [y/n] ")
	if err != nil {
		return err
	}

	if addUpperCase == `y` {
		passwordBytes += upperCaseBytes
	}

	addDigits, err := readString("Digits? [y/n] ")
	if err != nil {
		return err
	}

	if addDigits == `y` {
		passwordBytes += digitBytes
	}

	addSpecial, err := readString("Special characters? [y/n] ")
	if err != nil {
		return err
	}

	if addSpecial == `y` {
		passwordBytes += specialBytes
	}

	password, err := generatePassword(passwordBytes, passwordLength)
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(password); err != nil {
		return err
	}

	fmt.Println("Copied password to clipboard")

	return nil
}

func generatePassword(baseBytes string, length int) (string, error) {
	var password string

	for range length {
		val, err := rand.Int(rand.Reader, big.NewInt(int64(len(baseBytes))))
		if err != nil {
			return "", err
		}

		password += string(baseBytes[val.Int64()])
	}

	return password, nil
}
