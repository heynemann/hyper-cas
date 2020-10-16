package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vtex/hyper-cas/synchronizer"
)

var syncLabel string
var syncURL string
var syncJson bool

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
		result, err := s.Run(syncLabel)
		if err != nil {
			panic(err)
		}
		if syncJson {
			res, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(res))
		} else {
			printResult(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().StringVarP(&syncLabel, "label", "l", "", "Label to apply to this new distribution")
	syncCmd.Flags().StringVarP(&syncURL, "api-url", "a", "http://localhost:2485/", "Hyper-CAS API URL")
	syncCmd.Flags().BoolVarP(&syncJson, "json", "j", false, "Whether to output JSON serialization")
}

func printResult(result map[string]interface{}) {
	for _, file := range result["files"].([]map[string]interface{}) {
		isUpToDate := file["upToDate"].(bool)
		path := file["path"].(string)
		if isUpToDate {
			fmt.Printf("* %s - Already up-to-date.\n", path)
		} else {
			fmt.Printf("* %s - Updated (hash: %s).\n", path, file["hash"].(string))
		}
	}

	distro := result["distro"].(map[string]interface{})
	fmt.Printf("* Distro %s is up-to-date.\n", distro["hash"].(string))

	label := result["label"].(map[string]interface{})
	fmt.Printf("* Updated label %s => %s.\n", label["label"].(string), label["hash"].(string))
}
