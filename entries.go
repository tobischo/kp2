package main

import (
	"fmt"
	"strconv"
	"strings"

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

	entries := searchEntries(selectors, g)

	for i, entry := range entries {
		fmt.Printf("%3d %s\n", i, entry)
	}

	fmt.Print("Selection: ")
	var response string

	_, err := fmt.Scanln(&response)
	if err != nil {
		return nil, err
	}

	fmt.Println()

	index, err := strconv.Atoi(response)
	if err != nil {
		return nil, err
	}

	return readEntry(strings.Split(entries[index], "/"), g)

	return nil, fmt.Errorf("Failed to locate entry at selector")
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
