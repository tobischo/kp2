package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var cmdBrowse = &cobra.Command{
		Use:   "browse",
		Short: "interactive browsing mode",
		Long:  "browse is for interactively browsing the Keepass2 file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("browse")
		},
	}

	var cmdCopy = &cobra.Command{
		Use:   "copy [selector]",
		Short: "copies the password into the clipboard",
		Long:  `copy is for selecting the entry and copying the password into the clipboard`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("copy")
		},
	}

	var cmdRemove = &cobra.Command{
		Use:   "remove [selector]",
		Short: "removes an entry from the keepass file",
		Long:  `remove takes an entry out of the Keepass file. It asks for confirmation before persisting the file`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("remove")
		},
	}

	var cmdMove = &cobra.Command{
		Use:   "move [sourceSelector] [targetSelector]",
		Short: "moves an entry within the file",
		Long:  `move takes an entry from the position given with [sourceSelector] and moves it to the group given at [targetSelector]`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("move")
		},
	}

	var cmdCreate = &cobra.Command{
		Use:   "create [selector]",
		Short: "creates a new entry at the given location",
		Long:  `create builds a new entry at the given location and asks for the information required`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create")
		},
	}

	var cmdSearch = &cobra.Command{
		Use:   "search [selector]",
		Short: "looks through groups and entries",
		Long:  `search returns a list of entry selectors matching the given selector pattern`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("search")
		},
	}

	var cmdGeneratePassword = &cobra.Command{
		Use:   "generate [selector]",
		Short: "generates a new password",
		Long:  `generate builds a new password for the selected entry and copies it into the clipboard. It has to be accepted before it is persisted`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("generate password")
		},
	}

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "shows the version",
		Long:  "version shows the version of this tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version")
		}}

	var rootCmd = &cobra.Command{
		Use:   "kp2",
		Short: "tool to access Keepass2 files form the command line",
	}
	rootCmd.AddCommand(
		cmdBrowse, cmdCopy, cmdRemove, cmdMove, cmdCreate, cmdSearch,
		cmdGeneratePassword, cmdVersion,
	)
	rootCmd.Execute()
}
