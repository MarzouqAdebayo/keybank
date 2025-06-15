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
}

var newRemoteProfile = &cobra.Command{
	Use:   "new",
	Short: "New remote profile",
	Long:  "Add new remote profile",
	Run:   core.AddNewRemoteProfile,
}

// ssh
var sshCmd = &cobra.Command{
	Use:                "ssh",
	Short:              "Connect with profile",
	Long:               "Connect to a remote server via ssh using saved profile",
	DisableFlagParsing: true,
	Run:                core.SSHIntoProfile,
}

func init() {
	rootCmd.AddCommand(keysCmd)
	rootCmd.AddCommand(remoteCmd)

	rootCmd.AddCommand(sshCmd)

	// Remote subcommands
	newRemoteProfile.Flags().StringP("tag", "t", "", "Add a unique tag to the profile")
	newRemoteProfile.Flags().StringP("host", "r", "", "Remote profile host")
	newRemoteProfile.Flags().StringP("user", "u", "", "Remote profile user")
	newRemoteProfile.Flags().StringP("port", "p", "", "Remote profile port")
	newRemoteProfile.Flags().StringP("ssh_key", "s", "", "Remote profile ssh key file path")
	newRemoteProfile.MarkFlagsRequiredTogether("tag", "host", "user", "ssh_key")
	remoteCmd.AddCommand(newRemoteProfile)

	// keys subcommands
	keysCmd.AddCommand(listCmd)

	// list subcommands
	listCmd.AddCommand(listSSHCmd)
	listCmd.AddCommand(listGPGCmd)

}
