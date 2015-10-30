package main

import (
	"github.com/spf13/cobra"
	"github.com/tobischo/gokeepasslib"
)

const (
	app        = "pk2"
	timeFormat = "2006-01-02 15:04:05 -0700"
)

var usePassword bool
var keyFile string
var filePath string
var groupFlag bool

var changed bool = false

var db *gokeepasslib.Database

func main() {

	var cmdAdd = &cobra.Command{
		Use:     "add [selector]",
		Short:   "adds a new entry at the given location",
		Long:    `add builds a new entry at the given location and asks for the information required`,
		PreRunE: loadDatabaseCmd,
		RunE:    addCmd,
	}
	cmdAdd.Flags().BoolVarP(&groupFlag, "group", "g", false, "if set to true adds a group instead of an entry")

	// var cmdBrowse = &cobra.Command{
	// 	Use:   "browse",
	// 	Short: "interactive browsing mode",
	// 	Long:  "browse is for interactively browsing the Keepass2 file",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("browse")
	// 	},
	// }

	var cmdCopy = &cobra.Command{
		Use:     "copy [selector]",
		Short:   "copies the password into the clipboard",
		Long:    `copy is for selecting the entry and copying the password into the clipboard`,
		PreRunE: loadDatabaseCmd,
		RunE:    copyCmd,
	}

	// var cmdGeneratePassword = &cobra.Command{
	// 	Use:   "generate [selector]",
	// 	Short: "generates a new password",
	// 	Long:  `generate builds a new password for the selected entry and copies it into the clipboard. It has to be accepted before it is persisted`,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("generate password")
	// 	},
	// }

	var cmdInfo = &cobra.Command{
		Use:     "info [selector]",
		Short:   "shows the information for an entry",
		Long:    `info is for listing all relevant information for a command (except the password)`,
		PreRunE: loadDatabaseCmd,
		RunE:    infoCmd,
	}

	var cmdInit = &cobra.Command{
		Use:   "init [no options!]",
		Short: "initializes a new kdbx file",
		Long:  `init is for creating a new basic keepass file at the given location. It will fail if a file already exists`,
		RunE:  initCmd,
	}

	// var cmdMove = &cobra.Command{
	// 	Use:   "move [sourceSelector] [targetSelector]",
	// 	Short: "moves an entry within the file",
	// 	Long:  `move takes an entry from the position given with [sourceSelector] and moves it to the group given at [targetSelector]`,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("move")
	// 	},
	// }

	// var cmdRemove = &cobra.Command{
	// 	Use:   "remove [selector]",
	// 	Short: "removes an entry from the keepass file",
	// 	Long:  `remove takes an entry out of the Keepass file. It asks for confirmation before persisting the file`,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("remove")
	// 	},
	// }

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "shows the version",
		Long:  "version shows the version of this tool",
		Run:   versionCmd,
	}

	var rootCmd = &cobra.Command{
		Use:               "kp2",
		Short:             "tool to access Keepass2 files form the command line",
		PersistentPreRun:  paramSetupCmd,
		PersistentPostRun: persistDatabaseIfChanged,
	}

	rootCmd.PersistentFlags().BoolVarP(
		&usePassword, "password", "p", false,
		"true by default, false when using key param, add when using key and password authentication",
	)
	rootCmd.PersistentFlags().StringVarP(
		&keyFile, "key", "k", "",
		"path to the key file to use for auth",
	)
	rootCmd.PersistentFlags().StringVarP(
		&filePath, "file", "f", "",
		"Keepass2 file to be loaded, setting KP2FILE allows ommiting this flag",
	)

	rootCmd.AddCommand(
		cmdAdd,
		// cmdBrowse,
		cmdCopy,
		// cmdGeneratePassword,
		cmdInfo,
		cmdInit,
		// cmdMove,
		// cmdRemove,
		cmdVersion,
	)
	rootCmd.Execute()
}
