package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib"
)

func addCmd(cmd *cobra.Command, args []string) error {
	selectors := strings.Split(strings.Join(args, " "), "/")

	group, err := readGroup(selectors, &db.Content.Root.Groups[0])
	if err != nil {
		return err
	}

	if groupFlag {
		groupName, err := readString("Group Name: ")
		if err != nil {
			return err
		}

		for _, g := range group.Groups {
			if g.Name == groupName {
				return fmt.Errorf("Group Name must be unique within a parent group")
			}
		}

		newGroup := gokeepasslib.NewGroup()
		newGroup.Name = groupName
		group.Groups = append(group.Groups, newGroup)
		changed = true
	} else {
		entryTitle, err := readString("Entry Title: ")
		if err != nil {
			return err
		}

		for _, e := range group.Entries {
			if e.GetTitle() == entryTitle {
				return fmt.Errorf("Entry Title must be unique within a parent group")
			}
		}

		var password string
		password, err = readPasswordWithConfirmation()
		if err != nil {
			return err
		}

		values := []gokeepasslib.ValueData{
			{Key: "Password", Value: gokeepasslib.V{Content: password, Protected: true}},
			{Key: "Title", Value: gokeepasslib.V{Content: entryTitle}},
		}

		newEntry := gokeepasslib.NewEntry()
		newEntry.Values = append(newEntry.Values, values...)
		group.Entries = append(group.Entries, newEntry)
		changed = true
	}

	return nil
}
