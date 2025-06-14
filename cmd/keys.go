/*
Copyright © 2025 Marzouq Adebayo marzouqaadebayo@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MarzouqAdebayo/keybank/internal/core"
)

// keysCmd represents the keys command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH & GPG keys",
	Long:  `Add, list, and remove SSH and GPG keys tracked by keybank`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys",
	Long:  "List all SSH or GPG keys managed by Keybank.",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add",
	Long:  "Add",
}

var addRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Add remote",
	Long:  "Add new remote profile",
	Run:   core.AddNewRemoteProfile,
}

var listSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "List SSH keys",
	Run:   core.RunListSSH,
}

var listGPGCmd = &cobra.Command{
	Use:   "gpg",
	Short: "List GPG keys",
	Run:   core.RunListGPG,
}

func init() {
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(addCmd)

	// Flags
	addRemoteCmd.Flags().StringP("tag", "t", "", "Add a unique tag to the profile")
	addRemoteCmd.Flags().StringP("host", "r", "", "Remote profile host")
	addRemoteCmd.Flags().StringP("user", "u", "", "Remote profile user")
	addRemoteCmd.Flags().StringP("port", "p", "", "Remote profile port")
	addRemoteCmd.Flags().StringP("ssh_key", "s", "", "Remote profile ssh key file path")
	addRemoteCmd.MarkFlagsRequiredTogether("tag", "host", "user", "port", "ssh_key")

	// keys subcommands
	keysCmd.AddCommand(listCmd)
	listCmd.AddCommand(listSSHCmd)
	listCmd.AddCommand(listGPGCmd)
	addCmd.AddCommand(addRemoteCmd)
}
