/*
Copyright © 2022 SiteWhere LLC - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/devicechain-io/dc-microservice/core"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var ms *core.Microservice

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devicechain",
	Short: "DeviceChain Microservice",
	Long:  `Starts the DeviceChain microservice`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return createAndStartMicroservice()
	},
}

// Create microservice and initialize/start it.
func createAndStartMicroservice() error {
	log.Info().Msg("Creating new microservice and running intialization/startup...")
	ms = core.NewMicroservice()
	err := ms.Initialize(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("unable to initialize microservice")
		return err
	}
	err = ms.Start(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("unable to start microservice")
		return err
	}
	ms.WaitForShutdown()
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devicechain.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".devicechain" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".devicechain")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
