## WTF

A blazing fast command-line tool that automatically detects errors in previously executed console commands, suggests corrections, and provides an option to re-execute the corrected commands.

### Features

- 🚀 **Blazing Fast**: Reads only the last command from history without scanning the entire file
- 🎯 **Smart Suggestions**: Provides intelligent corrections for common typos and mistakes
- 🔄 **Quick Re-execution**: Re-run corrected commands with a simple confirmation
- ⚙️ **Customizable**: Define your own correction rules via configuration
- 🛡️ **Safe**: Asks for confirmation before executing any corrected command
- 🔌 **Shell Support**: Works with bash, zsh, and other common shells

### Installation

```bash
go install github.com/nowayhecodes/wtf@latest
```

Or quickly install using the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/nowayhecodes/wtf/main/install.sh | sh
```

Or build from source:

```bash
git clone https://github.com/nowayhecodes/wtf.git
cd wtf
go build -o bin/wtf cmd/wtf/main.go
```

### Usage

After a failed command, simply run:

```bash
wtf
```

The tool will:
1. Detect the last failed command
2. Suggest a correction
3. Execute the corrected command with your confirmation

#### Examples

```bash
$ gti status
git: command not found

$ wtf
Did you mean: git status? [Y/n] y
# Executes: git status
```

```bash
$ grpe "pattern" file.txt
grpe: command not found

$ wtf
Did you mean: grep "pattern" file.txt? [Y/n] y
# Executes: grep "pattern" file.txt
```

### Configuration

Create a `.wtf.json` file in your home directory:

```json
{
    "customRules": {
        "gti": "git",
        "sl": "ls"
    },
    "shellType": "bash",
    "maxSuggestions": 3,
    "levenThreshold": 2
}
```

#### Configuration Options

- `customRules`: Map of custom corrections
- `shellType`: Your shell type (bash, zsh)
- `maxSuggestions`: Maximum number of suggestions to show
- `levenThreshold`: Maximum edit distance for suggestions

### Development

Requirements:
- Go 1.22 or newer
- Make or Task (for build automation)

Setup:

```bash
# Clone the repository
git clone https://github.com/nowayhecodes/wtf.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
task build
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Acknowledgments

- Inspired by various command-line tools that improve developer productivity

### Author

[@nowayhecodes](https://github.com/nowayhecodes)

### Support

If you like this project, please give it a ⭐️!