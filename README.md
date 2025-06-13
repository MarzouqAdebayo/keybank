````markdown
# keybank

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![CI](https://img.shields.io/github/actions/workflow/status/MarzouqAdebayo/keybank/ci.yml?branch=main) ![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

**keybank** is a Go-based CLI tool for managing your SSH and GPG keys, associating them with remote hosts, and launching SSH sessions using saved profiles.

## Features

- List and add SSH & GPG keys.
- Create, list, and remove remote connection profiles.
- Launch SSH sessions using saved profiles.

## Installation

```bash
go install github.com/yourusername/keybank@latest
```
````

## Usage

```bash
# List SSH keys
keybank keys list ssh

# Add a GPG key
keybank keys add gpg --id 0xABCD1234

# Create a remote profile
keybank remote add prod-db --host 10.0.1.4 --user ubuntu --ssh-key ~/.ssh/prod.pem

# Connect via SSH
keybank ssh prod-db
```

## Configuration

Configuration files and directories are stored under `$HOME/.keybank/`:

- **config.yaml** or **config.json** for global settings.
- **keys/ssh/** and **keys/gpg/** for indexed key metadata.
- **remotes/** for remote profile definitions.

## Development

To contribute or extend keybank locally:

```bash
git clone https://github.com/MarzouqAdebayo/keybank.git
cd keybank
make setup      # install dependencies (optional)
make test       # run unit tests
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on reporting issues and submitting pull requests.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
