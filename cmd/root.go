/*
Copyright © 2025 Marzouq Adebayo marzouqaadebayo@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appDir     string
	configFile string
)

const defaultConfigFile = ".keybank"

var rootCmd = &cobra.Command{
	Use:   "keybank",
	Short: "Manage SSH & GPG keys, remote profiles, and SSH sessions",
	Long: `Keybank is a CLI tool for securely managing your SSH and GPG keys,
associating them with remote hosts, and launching SSH sessions using
saved connection profiles.

Examples:
  # List all SSH keys
  keybank keys list ssh

  # Add a GPG key by fingerprint
  keybank keys add gpg --id ABCD1234

  # Create a remote profile for production
  keybank remote add prod-db --host example.com --user ubuntu --ssh-key ~/.ssh/prod.pem

  # Connect to production database via SSH
  keybank ssh prod-db
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Welcome to keybank, use keybank --help to get started")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultAppDir := filepath.Join(home, defaultConfigFile)

	// Global persistent flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.keybank/.keybank.yaml)")
	rootCmd.PersistentFlags().StringVar(&appDir, "app-dir", defaultAppDir, "where keybank stores all its data")

	viper.BindPFlag("appDir", rootCmd.PersistentFlags().Lookup("app-dir"))
	viper.SetDefault("appDir", defaultAppDir)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		defaultAppDir := filepath.Join(home, defaultConfigFile)

		viper.AddConfigPath(".")
		viper.AddConfigPath(defaultAppDir)
		ext := strings.ToLower(filepath.Ext(configFile))
		switch ext {
		case ".yaml", ".yml":
			viper.SetConfigType("yaml")
		case ".json":
			viper.SetConfigType("json")
		default:
		}
		viper.SetConfigName("keybank")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
