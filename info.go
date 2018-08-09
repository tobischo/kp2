package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func infoCmd(cmd *cobra.Command, args []string) error {
	selectors := strings.Split(strings.Join(args, " "), "/")

	entry, err := readEntry(selectors, &db.Content.Root.Groups[0])
	if err != nil {
		return err
	}

	fmt.Printf("Title:             %s\n", entry.GetTitle())
	fmt.Printf("Creation:          %s\n", time.Time(*entry.Times.CreationTime).Format(timeFormat))
	fmt.Printf("Last Modification: %s\n", time.Time(*entry.Times.LastModificationTime).Format(timeFormat))
	fmt.Printf("Last Access:       %s\n", time.Time(*entry.Times.LastAccessTime).Format(timeFormat))
	fmt.Printf("UserName:          %s\n", entry.GetContent("UserName"))
	fmt.Printf("URL:               %s\n", entry.GetContent("URL"))

	fmt.Printf("Notes:\n%s\n", entry.GetContent("Notes"))

	markAsAccessed(entry)

	return nil
}
