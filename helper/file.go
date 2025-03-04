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
	
	// Reorder dirs to put root directory at the end
	var rootIndex int
	var hasRootFiles bool
	for i, dir := range dirs {
		if dir == "" {
			rootIndex = i
			hasRootFiles = true
			break
		}
	}
	
	// If root files exist, move them to the end
	if hasRootFiles {
		dirs = append(append(dirs[:rootIndex], dirs[rootIndex+1:]...), "")
	}
	
	// Print root directory name
	fmt.Println(".")
	
	// Print files by directory
	for i, dir := range dirs {
		dirFiles := filesByDir[dir]
		isLastDir := i == len(dirs)-1
		
		if dir == "" {
			// Root directory files (now at the end)
			for j, file := range dirFiles {
				prefix := "├── "
				if j == len(dirFiles)-1 {
					prefix = "└── "
				}
				fmt.Println(prefix + file.File)
			}
		} else {
			// Calculate directory nesting level
			level := strings.Count(dir, string(os.PathSeparator)) + 1
			
			// Adjust indentation for the last directory group if no root files
			var indent string
			if isLastDir && !hasRootFiles && level > 1 {
				indent = strings.Repeat("│   ", level-2) + "    "
			} else {
				indent = strings.Repeat("│   ", level-1)
			}
			
			// Print directory name with proper indentation
			dirPrefix := "├── "
			if isLastDir && !hasRootFiles {
				dirPrefix = "└── "
			}
			
			// For directories after root level
			if level == 1 {
				fmt.Println(dirPrefix + dir + "/")
			} else {
				pathParts := strings.Split(dir, string(os.PathSeparator))
				fmt.Println(indent + dirPrefix + pathParts[len(pathParts)-1] + "/")
			}
			
			// Print files in this directory
			for j, file := range dirFiles {
				// Adjust file indentation for the last directory group if no root files
				var fileIndent string
				if isLastDir && !hasRootFiles {
					fileIndent = strings.Repeat("│   ", level-1) + "    "
				} else {
					fileIndent = strings.Repeat("│   ", level)
				}
				
				filePrefix := "├── "
				if j == len(dirFiles)-1 {
					filePrefix = "└── "
				}
				fmt.Println(fileIndent + filePrefix + file.File)
			}
		}
	}
}

// GetSlugifiedTargetPath calculates the target path with slugified directory names
func GetSlugifiedTargetPath(sourceBase, filePath, outputDir string, slugFunc func(string) string) string {
    // Normalize paths
    sourceBase = filepath.Clean(sourceBase)
    filePath = filepath.Clean(filePath)
    
    // Get the relative path from source base directory
    relPath, err := filepath.Rel(sourceBase, filepath.Dir(filePath))
    if err != nil {
        // Fallback if we can't get relative path
        return outputDir
    }
    
    if relPath == "." {
        return outputDir
    }
    
    // Split the path and slugify each directory name
    parts := strings.Split(relPath, string(os.PathSeparator))
    for i, part := range parts {
        parts[i] = slugFunc(part)
    }
    
    // Join with the correct separator
    slugifiedPath := strings.Join(parts, string(os.PathSeparator))
    
    // Combine output directory with slugified path
    return filepath.Join(outputDir, slugifiedPath)
}

// SortDirsByDepth sorts directories by depth, deepest first
func SortDirsByDepth(dirs []string) {
    // Sort directories by depth (deepest first)
    // to avoid renaming parent before child
    sort.Slice(dirs, func(i, j int) bool {
        // Count separators to determine depth
        depthI := strings.Count(dirs[i], string(os.PathSeparator))
        depthJ := strings.Count(dirs[j], string(os.PathSeparator))
        
        // If equal depth, sort alphabetically for deterministic ordering
        if depthI == depthJ {
            return dirs[i] > dirs[j] // Reverse alphabetical
        }
        
        return depthI > depthJ // Higher depth comes first
    })
}

// GetDirectories returns a list of unique directories from given files
// Excludes the base directory if provided
func GetDirectories(files []model.File, baseDir string) []string {
    dirMap := make(map[string]bool)
    
    for _, file := range files {
        dirPath := file.Folder
        if dirPath != baseDir { // Skip the base directory
            dirMap[dirPath] = true
        }
    }
    
    // Convert map to slice
    dirs := make([]string, 0, len(dirMap))
    for dir := range dirMap {
        dirs = append(dirs, dir)
    }
    
    return dirs
}

// MoveFilesByPath moves all file from source path to destination path
func MoveFilesByPath(src, dst string) error {
	// Check if source exists
	if !HasDir(src) {
		return fmt.Errorf("source directory does not exist: %s", src)
	}
	
	// Get all files recursively from source
	files, err := GetFiles(src, true)
	if err != nil {
		return err
	}
	
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	
	// Ensure paths end with separator
	if !strings.HasSuffix(src, string(os.PathSeparator)) {
		src += string(os.PathSeparator)
	}
	if !strings.HasSuffix(dst, string(os.PathSeparator)) {
		dst += string(os.PathSeparator)
	}
	
	// Move each file
	for _, file := range files {
		// Get the relative path
		relPath := strings.TrimPrefix(file.FullPath, src)
		destPath := filepath.Join(dst, relPath)
		
		// Ensure destination directory exists
		destDir := filepath.Dir(destPath)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}
		
		// Move the file
		if err := moveFile(file.FullPath, destPath); err != nil {
			return err
		}
	}
	
	return nil
}

// moveFile moves a single file from src to dst
func moveFile(src, dst string) error {
	// First try a simple rename
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	
	// If rename fails, try copy + delete approach
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	// Copy contents
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	
	// Close both files before removing source
	srcFile.Close()
	dstFile.Close()
	
	// Preserve file permissions
	srcInfo, err := os.Stat(src)
	if err == nil {
		os.Chmod(dst, srcInfo.Mode())
	}
	
	// Delete original file
	return os.Remove(src)
}