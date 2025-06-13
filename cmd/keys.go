/*
Copyright © 2025 Marzouq Adebayo marzouqaadebayo@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

var listSSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "List SSH keys",
	Run:   runListSSH,
}

var listGPGCmd = &cobra.Command{
	Use:   "gpg",
	Short: "List GPG keys",
	Run:   runListGPG,
}

func init() {
	rootCmd.AddCommand(keysCmd)

	// keys subcommands
	keysCmd.AddCommand(listCmd)
	listCmd.AddCommand(listSSHCmd)
	listCmd.AddCommand(listGPGCmd)
}

func runListSSH(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "▸ unable to find home directory: %v\n", err)
		os.Exit(1)
	}
	sshDir := filepath.Join(home, ".ssh")
	files, err := os.ReadDir(sshDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "▸ error scanning ~/.ssh: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SSH keys found:")
	for _, file := range files {
		filename := file.Name()
		base := filepath.Base(filename)
		// skipping pub keys and known non-key files
		if file.IsDir() || strings.HasSuffix(base, ".pub") || strings.HasPrefix(base, "known_hosts") || strings.HasPrefix(base, "config") || strings.HasSuffix(base, ".gitconfig") {
			continue
		}
		fmt.Printf("  • %s\n", base)
	}
}

func runListGPG(cmd *cobra.Command, args []string) {
	fmt.Println("GPG keys:")
	g := exec.Command("gpg", "--list-secret-keys", "--keyid-format=long")
	// g.Stdout = os.Stdout
	// g.Stderr = os.Stderr
	output, err := g.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "▸ error invoking gpg: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		fmt.Printf("%d. -- %s\n", i+1, line)
	}
}
