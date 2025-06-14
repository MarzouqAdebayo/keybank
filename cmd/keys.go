/*
Copyright © 2025 Marzouq Adebayo marzouqaadebayo@gmail.com
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MarzouqAdebayo/keybank/internal/core"
)

// keysCmd
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH & GPG keys",
	Long:  `Add, list, and remove SSH and GPG keys tracked by keybank`,
}

// listCmd
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys",
	Long:  "List all SSH or GPG keys managed by Keybank.",
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

// remoteCmd
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Add remote",
	Long:  "Add new remote profile",
	Run:   core.AddNewRemoteProfile,
}

var newRemoteProfile = &cobra.Command{
	Use:   "new",
	Short: "New remote profile",
	Long:  "Add new remote profile",
	Run:   core.AddNewRemoteProfile,
}

func init() {
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(remoteCmd)

	// Flags
	remoteCmd.Flags().StringP("tag", "t", "", "Add a unique tag to the profile")
	remoteCmd.Flags().StringP("host", "r", "", "Remote profile host")
	remoteCmd.Flags().StringP("user", "u", "", "Remote profile user")
	remoteCmd.Flags().StringP("port", "p", "", "Remote profile port")
	remoteCmd.Flags().StringP("ssh_key", "s", "", "Remote profile ssh key file path")
	remoteCmd.MarkFlagsRequiredTogether("tag", "host", "user", "port", "ssh_key")

	// keys subcommands
	keysCmd.AddCommand(listCmd)

	// list subcommands
	listCmd.AddCommand(listSSHCmd)
	listCmd.AddCommand(listGPGCmd)

	// Remote subcommands
	remoteCmd.AddCommand(newRemoteProfile)
}
