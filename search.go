package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func searchCmd(cmd *cobra.Command, args []string) {
	entries := listEntries(&db.Content.Root.Groups[0])

	searchString := strings.ToLower(strings.Join(args, " "))

	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry), searchString) {
			fmt.Println(entry)
		}
	}
}
