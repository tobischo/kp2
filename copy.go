package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func copyCmd(cmd *cobra.Command, args []string) error {
	selectors := strings.Split(strings.Join(args, " "), "/")

	entry, err := readEntry(selectors, &db.Content.Root.Groups[0])
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(entry.GetPassword()); err != nil {
		return err
	}

	markAsAccessed(entry)

	fmt.Printf("UserName: %s\n", entry.GetContent("UserName"))

	return nil
}
