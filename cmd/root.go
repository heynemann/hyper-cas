package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

var rootDebug bool
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hyper-cas",
	Short: "hyper-cas is a CAS server with distributions that's really fast",
	Long: `hyper-cas is a content-addressable storage that allows users to
upload content using the hash of the content as storage key and create
distributions of that content that can be served.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.LogError("Failed to run hyper-cas CLI.", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hyper-cas.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&rootDebug, "debug", "d", false, "Run hyper-cas in debug mode")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if rootDebug {
		utils.SetDebug()
	}

	viper.SetConfigType("yaml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("hyper-cas.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		utils.LogInfo("Successfully loaded config file.", zap.String("configPath", viper.ConfigFileUsed()))
	}
}
