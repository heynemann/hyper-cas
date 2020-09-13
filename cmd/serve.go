/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/heynemann/hyper-cas/cache"
	"github.com/heynemann/hyper-cas/hash"
	"github.com/heynemann/hyper-cas/serve"
	"github.com/heynemann/hyper-cas/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var servePort int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "hyper-cas server for storage API",
	Long: `hyper-cas serve handles all requests to store either data or
distributions.`,
	Run: func(cmd *cobra.Command, args []string) {
		hasherType := hash.SHA256
		cacheType := cache.LRU
		storageType := storage.Memory
		if viper.GetString("storage.type") == "file" {
			storageType = storage.FileSystem
		}
		app, err := serve.NewApp(servePort, hasherType, storageType, cacheType)
		if err != nil {
			fmt.Printf("Starting hyper-cas storage API failed with: %v", err)
			os.Exit(1)
		}
		app.ListenAndServe()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 2485, "Port to run hyper-cas API in")
}
