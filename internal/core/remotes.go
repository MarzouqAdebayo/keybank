package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/MarzouqAdebayo/keybank/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	DefaultUser        = "root"
	DefaultPort        = "22"
	ProfilesFileName   = "profiles.yaml"
	ConfigViperDirFlag = "appDir"
	ConfigFileName     = ".keybank.yaml"
)

type RemoteProfile struct {
	ID     int    `yaml:"id" json:"id"`
	Tag    string `yaml:"tag" json:"tag"`
	Host   string `yaml:"host" json:"host"`
	User   string `yaml:"user" json:"user"`
	Port   string `yaml:"port" json:"port"`
	SSHKey string `yaml:"ssh_key" json:"ssh_key"`
	GPGKey string `yaml:"gpg_key,omitempty" json:"gpg_key,omitempty"`
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

func ParseTagforSSH(args []string) (string, []string, error) {
	var tag string

	flagSet := pflag.NewFlagSet("ssh", pflag.ContinueOnError)
	flagSet.StringP("tag", "t", "", "Remote profile tag")
	flagSet.SetOutput(io.Discard)

	_ = flagSet.Parse(args)

	if len(tag) == 0 {
		return "", nil, errors.New("tag flag is required")
	}

	remainingArgs := []string{}
	index := 0
	for index+1 < len(args) {
		arg := args[index]
		nextArg := args[index+1]
		if strings.HasPrefix(arg, "---") {
			return "", nil, errors.New("Invalid flags passed")
		}

		if arg == "-tag" || arg == "--tag" {
			if strings.HasPrefix(nextArg, "-") {
				return "", nil, errors.New("No arguments passed for reqired flag: tag")
			} else {
				index += 2
			}
		} else {
			remainingArgs = append(remainingArgs, arg)
			index += 1
		}
	}
	remainingArgs = append(remainingArgs, args[index])

	return tag, remainingArgs, nil
}

func SSHIntoProfile(cmd *cobra.Command, args []string) {
	template := "Usage:\n keybank ssh [flags]\n\nFlags:\n -tag remote profile tag\n [flag] underlying ssh flags\n"

	tag, remainingArgs, err := ParseTagforSSH(args)

	if err != nil {
		cmd.Println(template)
		return
	}

	configDir := viper.GetString(ConfigViperDirFlag)
	profileFilePath := filepath.Join(configDir, ProfilesFileName)

	data, err := os.ReadFile(profileFilePath)
	if err != nil {
		cmd.Println("Error failed here 111")
		cmd.PrintErrln(err.Error())
		return
	}

	savedProfiles := &Config{}
	err = savedProfiles.Unmarshal(data)
	if err != nil {
		cmd.PrintErrln(err.Error())
		return
	}

	var profile *RemoteProfile
	for _, savedProfile := range savedProfiles.Profiles {
		if savedProfile.Tag == tag {
			profile = &savedProfile
			break
		}
	}

	if profile == nil {
		cmd.Printf("Profile with tag '%s' does not exist", tag)
		return
	}

	sshBin, err := exec.LookPath("ssh")
	if err != nil {
		cmd.PrintErrf("ssh not found in PATH: %v\n", err)
		return
	}
	execArgs := []string{
		sshBin,
		"%s@%s", profile.User, profile.Host,
		"-i", profile.SSHKey,
		"-p", profile.Port,
	}
	execArgs = append(args, remainingArgs...)

	env := os.Environ()
	if err := syscall.Exec(sshBin, execArgs, env); err != nil {
		cmd.PrintErrf("Failed to exec ssh: %v\n", err)
	}
}

func AddNewRemoteProfile(cmd *cobra.Command, args []string) {
	tag, _ := cmd.Flags().GetString("tag")
	host, _ := cmd.Flags().GetString("host")
	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetString("port")
	sshKey, _ := cmd.Flags().GetString("ssh_key")

	for name, val := range map[string]string{
		"tag": tag, "host": host, "user": user, "ssh-key": sshKey,
	} {
		if val == "" {
			cmd.PrintErrf("Error: --%s is required\n", name)
			return
		}
	}

	newProfile := RemoteProfile{
		Tag:    tag,
		Host:   host,
		User:   user,
		Port:   DefaultPort,
		SSHKey: sshKey,
	}
	if port != "" {
		newProfile.Port = port
	}

	base := viper.GetString(ConfigViperDirFlag)
	if err := utils.CreateDir(base); err != nil {
		cmd.PrintErrf("Error creating config dir %q: %v\n", base, err)
		return
	}
	cfgPath := filepath.Join(base, ProfilesFileName)

	profileFilePath := filepath.Join(cfgPath, "profiles.yaml")

	cfg := &Config{}
	if exists, err := utils.FileExists(profileFilePath); err != nil {
		cmd.PrintErrf("Error checking profiles file: %v\n", err)
		return
	} else if exists {
		data, err := os.ReadFile(profileFilePath)
		if err != nil {
			cmd.PrintErrf("Error reading config: %v\n", err)
			return
		}
		if err = cfg.Unmarshal(data); err != nil {
			cmd.PrintErrf("Error parsing config: %v\n", err)
			return
		}
	}

	for _, profile := range cfg.Profiles {
		if profile.Tag == newProfile.Tag {
			cmd.PrintErrln("You need to provide a unique tag")
			return
		}
	}

	if exists, err := utils.FileExists(newProfile.SSHKey); err != nil {
		cmd.PrintErrln(err.Error())
		return
	} else if !exists {
		cmd.PrintErrf("Error: SSH key %q does not exist\n", newProfile.SSHKey)
		return
	}

	if len(cfg.Profiles) > 0 {
		newProfile.ID = cfg.Profiles[len(cfg.Profiles)-1].ID + 1
	}
	cfg.Profiles = append(cfg.Profiles, newProfile)

	out, err := cfg.Marshal()
	if err != nil {
		cmd.PrintErrf("Error serializing config: %v\n", err)
		return
	}

	if err = os.WriteFile(profileFilePath, out, 0o600); err != nil {
		cmd.PrintErrf("Error writing config: %v\n", err)
		return
	}

	cmd.Println(fmt.Sprintf("✅ New profile %q added\n", newProfile.Tag))
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
