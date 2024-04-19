package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
	"github.com/tobischo/gokeepasslib/v3/wrappers"
)

func markAsAccessed(entry *gokeepasslib.Entry) {
	access := wrappers.Now()
	entry.Times.LastAccessTime = &access
	changed = true
}

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

func readEntry(selection string, g *gokeepasslib.Group) (*gokeepasslib.Entry, error) {
	for i, entry := range g.Entries {
		if entry.GetTitle() == selection {
			return &g.Entries[i], nil
		}
	}

	selectors := strings.Split(selection, "/")

	if len(selectors) == 1 {
		for i, entry := range g.Entries {
			if entry.GetTitle() == selectors[0] {
				return &g.Entries[i], nil
			}
		}
	} else {
		for i, group := range g.Groups {
			if group.Name == selectors[0] {
				return readEntry(strings.Join(selectors[1:], "/"), &g.Groups[i])
			}
		}
	}

	entries := searchEntries(selectors, g)
	if len(entries) < 1 {
		return nil, fmt.Errorf("No entry found")
	}

	if len(entries) == 1 {
		return readEntry(entries[0], g)
	}

	for i, entry := range entries {
		fmt.Printf("%3d %s\n", i, entry)
	}

	selection, err := readString("Selection: ")
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(selection)
	if err != nil {
		return nil, err
	}

	return readEntry(entries[index], g)
}

func searchEntries(selectors []string, g *gokeepasslib.Group) []string {
	entries := listEntries(g)

	selector := strings.ToLower(strings.Join(selectors, "/"))

	var selectedEntries = make([]string, 0)

	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry), selector) {
			selectedEntries = append(selectedEntries, entry)
		}
	}

	return selectedEntries
}
