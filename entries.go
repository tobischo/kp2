package main

import (
	"fmt"

	"github.com/tobischo/gokeepasslib"
)

func listEntries(g *gokeepasslib.Group) []string {
	var entries = make([]string, 0)

	for _, entry := range g.Entries {
		entries = append(entries, entry.GetTitle())
	}

	for _, group := range g.Groups {
		subEntries := listEntries(&group)
		for i, val := range subEntries {
			subEntries[i] = fmt.Sprintf("%s/%s", group.Name, val)
		}

		entries = append(entries, subEntries...)
	}

	return entries
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
