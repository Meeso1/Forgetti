# Forgetti

A tool that encrypts your data, and sometimes decrypts it too.

Forgetti lets you encrypt any data using a password and a server-based key. The data can then be decrypted within a configured time limit - after configured time runs out, it is inaccessible forever!
No server trust is needed - your data or your password is never sent the server. 

## Installation

```bash
# Clone the repository
git clone https://github.com/Meeso1/Forgetti.git
cd Forgetti

# Build everything (CLI tool and server binary)
make build

# Run the server
./bin/forgetti-server

# Run the CLI
./bin/forgetti-cli --version
```

## Usage

```bash
# Coming soon...
```

## Development

```bash
# Get dependencies
make deps

# Run tests
make test
```

## Building with Version Information

To build with version information:

```bash
VERSION=1.0.0 && make build
``` 