package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func persistDatabaseIfChanged(cmd *cobra.Command, args []string) {
	if changed {
		if err := db.LockProtectedEntries(); err != nil {
			fmt.Printf("Failed to lock entries: %s\n", err)
			return
		}

		/*		file, err := os.Create(filePath)
				if err != nil {
					fmt.Printf("Failed to open file to persist changes: %s\n", err)
					return
				}
				defer file.Close()

				encoder := gokeepasslib.NewEncoder(file)
				err = encoder.Encode(db)
				if err != nil {
					fmt.Printf("Failed to encode database with: %s\n", err)
				}*/
	}
}
