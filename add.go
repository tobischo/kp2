package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

func addCmd(_ *cobra.Command, args []string) error {
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
				return errGroupNameNotUnique
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

		entryUserName, err := readString("Entry UserName: ")
		if err != nil {
			return err
		}

		entryURL, err := readString("Entry URL: ")
		if err != nil {
			return err
		}

		for _, e := range group.Entries {
			if e.GetTitle() == entryTitle {
				return errEntryTitleNotUnique
			}
		}

		var password string
		password, err = readPasswordWithConfirmation()
		if err != nil {
			return err
		}

		values := []gokeepasslib.ValueData{
			{Key: "Notes", Value: gokeepasslib.V{Content: "Notes"}},
			{
				Key: "Password",
				Value: gokeepasslib.V{
					Content:   password,
					Protected: wrappers.NewBoolWrapper(true),
				},
			},
			{Key: "Title", Value: gokeepasslib.V{Content: entryTitle}},
			{Key: "URL", Value: gokeepasslib.V{Content: entryURL}},
			{Key: "UserName", Value: gokeepasslib.V{Content: entryUserName}},
		}

		newEntry := gokeepasslib.NewEntry()
		newEntry.Values = append(newEntry.Values, values...)
		group.Entries = append(group.Entries, newEntry)
		changed = true
	}

	return nil
}
