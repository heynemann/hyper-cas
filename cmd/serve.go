package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vtex/hyper-cas/serve"
	"github.com/vtex/hyper-cas/storage"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
)

var servePort int
var profile bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "hyper-cas server for storage API",
	Long: `hyper-cas serve handles all requests to store either data or
distributions.`,
	Run: func(cmd *cobra.Command, args []string) {
		storageType := storage.FileSystem
		app, err := serve.NewApp(servePort, storageType)
		app.EnableProfileRoutes(profile)
		if err != nil {
			utils.LogError(
				"Starting hyper-cas storage API failed",
				zap.Error(err),
			)
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
	serveCmd.Flags().BoolVar(&profile, "profile", false, "Enable pprof debug routes")
}
