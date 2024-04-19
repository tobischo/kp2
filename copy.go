package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func copyCmd(cmd *cobra.Command, args []string) error {
	entry, err := readEntry(strings.Join(args, " "), &db.Content.Root.Groups[0])
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(entry.GetPassword()); err != nil {
		return err
	}

	markAsAccessed(entry)

	fmt.Printf("URL: %s\n", entry.GetContent("URL"))
	fmt.Printf("UserName: %s\n", entry.GetContent("UserName"))

	return nil
}
