package core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MarzouqAdebayo/keybank/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const DefaultUser = "root"
const DefaultPort = "22"

type RemoteProfile struct {
	ID     int    `yaml:"id" json:"id"`
	Tag    string `yaml:"tag" json:"tag"`
	Host   string `yaml:"host" json:"host"`
	User   string `yaml:"user" json:"user"`
	Port   string `yaml:"port" json:"port"`
	SSHKey string `yaml:"ssh_key" json:"ssh_key"`
	GPGKey string `yaml:"gpg_key,omitempty" json:"gpg_key,omitempty"`
}

func newRemoteProfile() RemoteProfile {
	return RemoteProfile{
		User: DefaultUser,
		Port: DefaultPort,
	}
}

func (rp RemoteProfile) Marshal() ([]byte, error) {
	return yaml.Marshal(rp)
}

func (rp RemoteProfile) Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, rp)
}

type Config struct {
	Profiles []RemoteProfile `yaml:"profiles"`
}

func (cfg *Config) Marshal() ([]byte, error) {
	return yaml.Marshal(cfg)
}

func (cfg *Config) Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, cfg)
}

func validateAddNewRemoteProfileFlags(cmd *cobra.Command, flag string) (string, error) {
	tag, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", err
	}
	if len(tag) == 0 {
		return "", errors.New(fmt.Sprintf("%s cannot be empty", flag))
	}
	return tag, nil
}

func AddNewRemoteProfile(cmd *cobra.Command, args []string) {
	tag, err := validateAddNewRemoteProfileFlags(cmd, "tag")
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}
	host, err := validateAddNewRemoteProfileFlags(cmd, "host")
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}
	user, err := validateAddNewRemoteProfileFlags(cmd, "user")
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}
	port, err := validateAddNewRemoteProfileFlags(cmd, "port")
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}
	ssh_key, err := validateAddNewRemoteProfileFlags(cmd, "ssh_key")
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}

	newProfile := newRemoteProfile()
	newProfile.Tag = tag
	newProfile.User = user
	newProfile.Host = host
	newProfile.Port = port
	newProfile.SSHKey = ssh_key

	configDir := viper.GetString("appDir")
	exists, err := utils.DirExists(configDir)
	if !exists && err != nil {
		cmd.PrintErrln(err.Error())
		return
	}

	if !exists {
		if err := utils.CreateDir(configDir); err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
	}

	profileFilePath := filepath.Join(configDir, "profiles.yaml")

	exists, err = utils.FileExists(profileFilePath)
	if !exists && err != nil {
		cmd.PrintErrln(err.Error())
		return
	}

	var data []byte
	if exists {
		data, err = os.ReadFile(profileFilePath)
		if err != nil {
			cmd.Println("Error failed here 111")
			cmd.PrintErrln(err.Error())
			return
		}
	}

	prevCfg := &Config{}
	err = prevCfg.Unmarshal(data)
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}

	for _, profile := range prevCfg.Profiles {
		if profile.Tag == newProfile.Tag {
			cmd.PrintErrln("You need to provide a unique tag")
			return
		}
	}

	if len(prevCfg.Profiles) == 0 {
		newProfile.ID = 0
	} else {
		newProfile.ID = prevCfg.Profiles[len(prevCfg.Profiles)-1].ID + 1
	}

	prevCfg.Profiles = append(prevCfg.Profiles, newProfile)
	b, err := prevCfg.Marshal()
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}

	err = os.WriteFile(profileFilePath, b, 0644)
	if err != nil {
		cmd.Println("Error failed here")
		cmd.PrintErrln(err.Error())
		return
	}
	cmd.Println(fmt.Sprintf("New profile with tag: %s added", newProfile.Tag))
}

func RunListSSH(cmd *cobra.Command, args []string) {
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

func RunListGPG(cmd *cobra.Command, args []string) {
	fmt.Println("GPG keys:")
	g := exec.Command("gpg", "--list-secret-keys", "--with-colons", "--keyid-format=long")
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
