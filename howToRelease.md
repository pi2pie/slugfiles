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
