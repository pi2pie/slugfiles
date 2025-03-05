# Slugfiles

A lightweight command-line utility for transforming file/directory names into clean, URL-friendly slugs.

>[!Note]
>
> The core functionality is provided by the [slug](https://github.com/gosimple/slug) package.
>
> The slug transformation process converts filenames to URL-friendly format by transliterating Unicode characters to ASCII (including converting non-ASCII characters to their English phonetic equivalents), converting spaces to hyphens, removing special characters, and ensuring lowercase output by default. Use the `--case-sensitive` flag (`-c`) to preserve uppercase letters in the output. More details see the [usage](#usage) section.


## Introduction

Slugfiles is a command-line tool written in Go that converts filenames and directory paths into URL-friendly slugs. It's perfect for bulk renaming files for web projects, content management systems, or any situation where you need normalized filenames. Slugfiles helps you:

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

### From Releases

Download the appropriate binary for your platform from the [releases page](https://github.com/pi2pie/slugfiles/releases).

## Uninstalling

To remove the Slugfiles binary from your system:

```bash
# Remove the binary 
go clean -i github.com/pi2pie/slugfiles

# Or manually delete the binary from your GOPATH/bin directory
rm $(go env GOPATH)/bin/slugfiles
```

## Usage

```bash
$ slugfiles

Rename files in a directory to user friendly slugs.

Usage:
  slugfiles [folder] [flags]
  slugfiles [command]

Available Commands:
  help        Help about any command
  rename      Rename files in a directory to user friendly slugs.

Flags:
  -c, --case-sensitive   Case sensitive renaming
  -d, --dry-run          Simulate renaming without making changes
  -h, --help             help for slugfiles
  -o, --output string    Output directory
  -r, --recursive        Process directories recursively
  -v, --version          version for slugfiles
```

You can use Slugfiles in two equivalent ways:

1. By providing a folder path directly:
   ```bash
   slugfiles /path/to/directory [flags]
   ```

2. By using the explicit `rename` command:
   ```bash
   slugfiles rename /path/to/directory [flags]
   ```

Both methods perform exactly the same operation. Choose whichever syntax you prefer.

### Recursive Mode

By default, Slugfiles only processes files in the specified directory without entering subdirectories. When you use the `--recursive` flag, it will process all files in the specified directory and all of its subdirectories.

### Output Directory

By default, Slugfiles renames files in the original directory. If you want to save the renamed files to a different directory, use the `--output` flag. The output directory must exist before running the command.

### Case Sensitivity

By default, Slugfiles converts filenames to lowercase. If you want to preserve the original case of the filenames, use the `--case-sensitive` flag.

### Dry Run Mode

The `--dry-run` flag allows you to preview the changes that would be made without actually renaming any files. This is useful for verifying what changes would occur before committing to them.

```bash
# Preview file renaming without making changes
slugfiles rename --dry-run /path/to/directory

# Preview recursive renaming
slugfiles rename --dry-run --recursive /path/to/directory
```

When in dry run mode, Slugfiles will display "[DRY RUN]" before each operation that would be performed.


### Some examples:

```bash
# Basic usage - rename files in a directory (only top-level files)
slugfiles /path/to/directory
# Or equivalently:
slugfiles rename /path/to/directory

# Rename recursively (includes files in all subdirectories)
slugfiles --recursive /path/to/directory
# Or equivalently:
slugfiles rename --recursive /path/to/directory

# Rename files and save to output directory
slugfiles /path/to/directory --output /path/to/output-directory
# Or equivalently:
slugfiles rename /path/to/directory --output /path/to/output-directory

# Rename files in all subdirectories and save to output directory
slugfiles --recursive /path/to/directory --output /path/to/output-directory
# Or equivalently:
slugfiles rename --recursive /path/to/directory --output /path/to/output-directory

# Rename files in a case-sensitive manner
slugfiles --case-sensitive /path/to/directory
# Or equivalently:
slugfiles rename --case-sensitive /path/to/directory

# Rename files in a case-sensitive manner and save to output directory
slugfiles --case-sensitive /path/to/directory --output /path/to/output-directory
# Or equivalently:
slugfiles rename --case-sensitive /path/to/directory --output /path/to/output-directory
```

## License

MIT