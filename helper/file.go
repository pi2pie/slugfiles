package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pi2pie/slugfiles/model"
)

// check isWindows or isUnix
func IsWindows() bool {
	return os.PathSeparator == '\\'
}

// GetSeparator : Get the separator for the OS
func GetSeparator() string {
	if IsWindows() {
		fmt.Println("OS: Windows")
		return "\\"
	}
	fmt.Println("OS: Unix")
	return "/"
}

// HasFile : Check if file exists in the current directory
func HasFile(filename string) bool {
	if info, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

// HasDir checks if a directory exists
func HasDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetFiles returns a list of files in a directory
// If recursive is true, it will recursively get files from subdirectories
func GetFiles(path string, recursive bool) ([]model.File, error) {
	files := []model.File{}

	// Check if path exists
	if !HasDir(path) {
		return files, fmt.Errorf("directory does not exist: %s", path)
	}

	// Ensure path ends with separator
	if !strings.HasSuffix(path, string(os.PathSeparator)) {
		path += string(os.PathSeparator)
	}

	if recursive {
		// Walk through all directories recursively
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			// Skip directories themselves
			if info.IsDir() {
				return nil
			}

			folder := filepath.Dir(filePath)
			if !strings.HasSuffix(folder, string(os.PathSeparator)) {
				folder += string(os.PathSeparator)
			}

			file := model.File{
				FullPath: filePath,
				Folder:   folder,
				File:     info.Name(),
				FileName: strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
				Ext:      filepath.Ext(info.Name()),
			}
			files = append(files, file)
			return nil
		})
		return files, err
	} else {
		// Non-recursive mode - only process files in the top directory
		fileInfos, err := os.ReadDir(path)
		if err != nil {
			return files, err
		}

		for _, fileInfo := range fileInfos {
			// Skip directories
			if fileInfo.IsDir() {
				continue
			}

			file := model.File{
				FullPath: path + fileInfo.Name(),
				Folder:   path,
				File:     fileInfo.Name(),
				FileName: strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name())),
				Ext:      filepath.Ext(fileInfo.Name()),
			}
			files = append(files, file)
		}
		return files, nil
	}
}

// GetNewFilePath :
func GetNewFilePath(file model.File, sourceFolder, targetFolder string) (string, error) {
	if file.FullPath == "" {
		return "", fmt.Errorf("original file is not valid")
	}

	subfolder := GetSubfolder(file, sourceFolder)

	return targetFolder + subfolder + file.File, nil
}

// GetSubfolder :
func GetSubfolder(file model.File, sourceFolder string) string {
	return strings.TrimPrefix(file.Folder, sourceFolder)
}

// MoveFile renames a file
func MoveFile(oldFile model.File, newFile model.File) error {
	return os.Rename(oldFile.FullPath, newFile.FullPath)
}

// CopyFile copies a file
func CopyFile(srcFile model.File, destFile model.File) error {
	src, err := os.Open(srcFile.FullPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destFile.Folder, 0755); err != nil {
		return err
	}

	dst, err := os.Create(destFile.FullPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// PrintFilesTree prints files in a tree structure grouped by directories
func PrintFilesTree(files []model.File, sourceDir string) {
	// Group files by directory for better organization
	filesByDir := make(map[string][]model.File)
	for _, file := range files {
		relDir, _ := filepath.Rel(sourceDir, file.Folder)
		if relDir == "." {
			relDir = "" // Root directory special case
		}
		filesByDir[relDir] = append(filesByDir[relDir], file)
	}
	
	// Sort directories for consistent display
	var dirs []string
	for dir := range filesByDir {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	
	// Print files by directory
	for i, dir := range dirs {
		if i > 0 {
			fmt.Println()
		}
		
		if dir == "" {
			// Root directory
			for j, file := range filesByDir[dir] {
				prefix := "├── "
				if j == len(filesByDir[dir])-1 {
					prefix = "└── "
				}
				fmt.Println(prefix + file.File)
			}
		} else {
			// Subdirectory
			fmt.Println(dir + "/")
			for j, file := range filesByDir[dir] {
				prefix := "  ├── "
				if j == len(filesByDir[dir])-1 {
					prefix = "  └── "
				}
				fmt.Println(prefix + file.File)
			}
		}
	}
}