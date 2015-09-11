package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib"
)

func copyCmd(cmd *cobra.Command, args []string) error {
	selectors := strings.Split(args[0], "/")

	if err := db.UnlockProtectedEntries(); err != nil {
		return err
	}

	entry, err := readEntry(selectors, &db.Content.Root.Groups[0])
	if err != nil {
		return err
	}

	fmt.Printf("%x", string(entry.GetPassword()))

	return nil
}

func readEntry(selectors []string, g *gokeepasslib.Group) (*gokeepasslib.Entry, error) {
	if len(selectors) == 1 {
		for _, entry := range g.Entries {
			if entry.GetTitle() == selectors[0] {
				return &entry, nil
			}
		}
	} else {
		for _, group := range g.Groups {
			if group.Name == selectors[0] {
				return readEntry(selectors[1:], &group)
			}
		}
	}

	return nil, fmt.Errorf("Failed to locate entry at selector")
}
