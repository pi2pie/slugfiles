# Slugfiles

A lightweight command-line utility for transforming file/directory names into clean, URL-friendly slugs.

## Introduction

Slugfiles is a command-line tool written in Go that converts filenames and directory paths into web-friendly slugs. It's perfect for bulk renaming files for web projects, content management systems, or any situation where you need normalized filenames. Slugfiles helps you:

- Remove special characters from filenames
- Convert spaces to hyphens
- Handle non-Latin characters appropriately
- Preserve file extensions
- Process files recursively (optional)
- Output to a different directory (optional)

## Installation

### Using Go

```bash
go install github.com/pi2pie/slugfiles@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/pi2pie/slugfiles.git
cd slugfiles

# Build and install
go install
```

After installation, verify it's working correctly:

```bash
slugfiles --version
```

### Uninstalling

To remove the Slugfiles binary from your system:

```bash
# Remove the binary 
go clean -i github.com/pi2pie/slugfiles

# Or manually delete the binary from your GOPATH/bin directory
rm $(go env GOPATH)/bin/slugfiles
```

### From Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/pi2pie/slugfiles/releases).

## Usage

```bash
# Basic usage - rename files in a directory (only top-level files)
slugfiles rename /path/to/directory

# Rename recursively (includes files in all subdirectories)
slugfiles rename --recursive /path/to/directory

# Rename files and save to output directory
slugfiles rename --output new-files /path/to/directory

# Rename files in all subdirectories and save to output directory
slugfiles rename --recursive --output new-files /path/to/directory
```

### Recursive Mode

By default, Slugfiles only processes files in the specified directory without entering subdirectories. When you use the `--recursive` flag, it will process all files in the specified directory and all of its subdirectories.

## Building and Releasing

### Setting Version

The version is defined in `cli/root.go`. To update the version:

1. Modify the Version field in RootCmd:
   ```go
   Version: "x.y.z",
   ```

2. Tag the repository with the same version:
   ```bash
   git tag -a vx.y.z -m "Release version x.y.z"
   git push origin vx.y.z
   ```

### Creating Releases

To build for multiple platforms:

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o dist/slugfiles-linux-amd64

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o dist/slugfiles-darwin-amd64

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o dist/slugfiles-windows-amd64.exe
```

For automated releases, consider setting up GitHub Actions with goreleaser.

## License

MIT