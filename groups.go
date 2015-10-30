package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tobischo/gokeepasslib"
)

func listGroups(g *gokeepasslib.Group) []string {
	var groups = make([]string, 0)

	groups = append(groups, g.Name)

	for _, group := range g.Groups {
		subGroups := listGroups(&group)

		for i, val := range subGroups {
			subGroups[i] = fmt.Sprintf("%s/%s", g.Name, val)
		}

		groups = append(groups, subGroups...)
	}

	return groups
}

func readGroup(selectors []string, g *gokeepasslib.Group) (*gokeepasslib.Group, error) {
	if len(selectors) == 1 {
		if g.Name == selectors[0] {
			return g, nil
		}
	} else {
		for i, group := range g.Groups {
			if group.Name == selectors[1] {
				return readGroup(selectors[1:], &g.Groups[i])
			}
		}
	}

	groups := searchGroups(selectors, g)
	if len(groups) < 1 {
		return nil, fmt.Errorf("No group found")
	}

	for i, group := range groups {
		fmt.Printf("%3d %s\n", i, group)
	}

	selection, err := readString("Selection: ")
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(selection)
	if err != nil {
		return nil, err
	}

	return readGroup(strings.Split(groups[index], "/"), g)

	return nil, fmt.Errorf("Failed to locate group at selector")
}

func searchGroups(selectors []string, g *gokeepasslib.Group) []string {
	groups := listGroups(g)

	selector := strings.ToLower(strings.Join(selectors, "/"))

	var selectedGroups = make([]string, 0)

	for _, group := range groups {
		if strings.Contains(strings.ToLower(group), selector) {
			selectedGroups = append(selectedGroups, group)
		}
	}

	return selectedGroups
}
