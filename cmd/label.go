package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vtex/hyper-cas/synchronizer"
)

var labelName string
var labelHash string
var labelRetries int
var labelURL string

// labelCmd represents the label command
var labelCmd = &cobra.Command{
	Use:   "set-label",
	Short: "Updates a label to a given distribution in hyper-cas",
	Long:  `set-label will set the specified label to the specified tree hash in hyper-cas`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		s := synchronizer.NewSync("", labelURL)
		retries := 0
		for i := 0; i <= labelRetries; i++ {
			hasDistro := s.HasDistro(labelHash)
			if !hasDistro {
				log.Fatalf("Distribution %s was not found\n", labelHash)
			}
			err = s.SetLabel(labelName, labelHash)
			if err == nil {
				break
			}
			retries += 1
		}
		if retries > labelRetries {
			log.Fatalf("Distribution %s could not be set to label %s: %v\n", labelHash, labelName, err)
		}
		fmt.Println("Label set successfully.")
	},
}

func init() {
	rootCmd.AddCommand(labelCmd)
	labelCmd.Flags().StringVarP(&labelURL, "api-url", "u", "http://localhost:2485/", "Hyper-CAS API URL")
	labelCmd.Flags().StringVarP(&labelName, "name", "n", "", "Label to set the hash of the distribution to")
	labelCmd.Flags().StringVarP(&labelHash, "hash", "a", "", "Distribution hash to set the label to")
	labelCmd.Flags().IntVarP(&labelRetries, "retries", "r", 3, "Number of times to retry setting the label")
}
