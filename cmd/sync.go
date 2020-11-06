package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vtex/hyper-cas/synchronizer"
)

var syncLabel string
var syncURL string
var syncJSON bool
var syncRetries int
var syncRequestRetries int
var syncMaxConcurrentRequests int
var syncHTTPTimeoutMs int
var syncDistroHTTPTimeoutMs int

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
		folder := strings.TrimRight(args[0], "/")
		if !filepath.IsAbs(folder) {
			folder, err = filepath.Abs(folder)
			if err != nil {
				panic(err)
			}
		}
		if !folderExists(folder) {
			panic(fmt.Sprintf("Folder %s does not exist!", folder))
		}
		s := synchronizer.NewSync(
			folder,
			syncURL,
			syncRequestRetries,
			syncMaxConcurrentRequests,
			syncHTTPTimeoutMs,
			syncDistroHTTPTimeoutMs,
		)
		var result map[string]interface{}
		retries := 0
		for i := 0; i <= syncRetries; i++ {
			result, err = s.Run(syncLabel)
			if err == nil {
				break
			}
			retries++
		}
		if result == nil {
			panic(fmt.Errorf("Failed to synchronize folder: %v", err))
		}
		result["retries"] = retries
		if syncJSON {
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
	syncCmd.Flags().StringVarP(&syncURL, "api-url", "u", "http://localhost:2485/", "Hyper-CAS API URL")
	syncCmd.Flags().BoolVarP(&syncJSON, "json", "j", false, "Whether to output JSON serialization")
	syncCmd.Flags().IntVarP(&syncRetries, "retries", "r", 0, "Number of times to retry the whole synchronizing process")
	syncCmd.Flags().IntVarP(&syncRequestRetries, "req-retries", "q", 0, "Number of times to retry each request to hyper-cas")
	syncCmd.Flags().IntVarP(&syncMaxConcurrentRequests, "max-concurrent", "m", 50, "Maximum number of concurrent requests to hyper-cas")
	syncCmd.Flags().IntVarP(&syncHTTPTimeoutMs, "timeout", "t", 5000, "Number of milliseconds to timeout per request to hyper-cas")
	syncCmd.Flags().IntVarP(&syncDistroHTTPTimeoutMs, "distro-timeout", "o", 300000, "Number of milliseconds to timeout when writing the distro to hyper-cas")
}

func printResult(result map[string]interface{}) {
	updatedFiles := result["files"].([]map[string]interface{})
	for _, file := range updatedFiles {
		path := file["path"].(string)
		if file["exists"] == true {
			fmt.Printf("* %s - Already exists (%vms).\n", path, file["duration"])
		} else {
			fmt.Printf("* %s - Updated (hash: %s - %vms).\n", path, file["hash"], file["duration"])
		}
	}

	distro := result["distro"].(map[string]interface{})
	fmt.Printf("* Distribution - Updated (hash: %s - %vms).\n", distro["hash"], distro["duration"])

	label := result["label"].(map[string]interface{})
	labelName := label["label"]
	if labelName != "" {
		fmt.Printf("* Updated label %s => %s.\n", labelName, label["hash"])
	}
	fmt.Printf(
		"Completed synchronizing %v files with %v retries in %vms.\n",
		len(updatedFiles),
		result["retries"],
		result["duration"],
	)
}
