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

The Forgetti CLI provides three main commands: `encrypt`, `decrypt`, and `metadata`.

### Encrypt a file

```bash
# Basic encryption 
./bin/forgetti-cli encrypt -i myfile.txt -o myfile.txt.forgetti

# Encrypt with custom expiration time (1 week)
./bin/forgetti-cli encrypt -i document.pdf -o document.pdf.forgetti -e 1w

# Encrypt with custom server
./bin/forgetti-cli encrypt -i secret.txt -o secret.txt.forgetti -s http://localhost:8080
```

### Decrypt a file

```bash
# Basic decryption
./bin/forgetti-cli decrypt -i myfile.txt.forgetti -o myfile_decrypted.txt

# Decrypt with custom server
./bin/forgetti-cli decrypt -i secret.txt.forgetti -o secret_restored.txt -s http://localhost:8080
```

### Read metadata from encrypted files

```bash
# View metadata of an encrypted file
./bin/forgetti-cli metadata -i myfile.txt.forgetti
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