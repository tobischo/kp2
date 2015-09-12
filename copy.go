package main

import (
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

	return nil
}
