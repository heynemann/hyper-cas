package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vtex/hyper-cas/synchronizer"
)

var syncLabel string
var syncURL string

func folderExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.Mode().IsDir()
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync a folder into a distribution in hyper-cas",
	Long:  `Sync will synchronize all files in a given folder into hyper-cas`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			panic("There must be only a single argument specifying path to sync.")
		}
		var err error
		folder := args[0]
		if !filepath.IsAbs(folder) {
			folder, err = filepath.Abs(folder)
			if err != nil {
				panic(err)
			}
		}
		if !folderExists(folder) {
			panic(fmt.Sprintf("Folder %s does not exist!", folder))
		}
		s := synchronizer.NewSync(folder, syncURL)
		_, err = s.Run(syncLabel)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&syncLabel, "label", "l", "", "Label to apply to this new distribution")
	syncCmd.Flags().StringVarP(&syncURL, "api-url", "a", "http://localhost:2485/", "Hyper-CAS API URL")
}
