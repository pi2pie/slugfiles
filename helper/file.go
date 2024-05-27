package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"pi2pie/slugfiles-rename/model"
)

// check isWindows or isUnix
func isWindows() bool {
	return os.PathSeparator == '\\'
}

// GetSeparator : Get the separator for the OS
func GetSeparator() string {
	if isWindows() {
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

// HasDir : Check if the dir is valid
func HasDir(dirname string) bool {
	if info, err := os.Stat(dirname); os.IsNotExist(err) {
		return false
	} else {
		return info.IsDir()
	}
}

// GetFiles :
func GetFiles(folder string, isRecursive bool) ([]model.File, error) {
	var files []model.File

	if isRecursive {
		err := filepath.Walk(folder, func(path string, info os.FileInfo, walkErr error) error {
			file, err := filepath.Abs(path)
			if err != nil {
				return nil
			}

			if !info.IsDir() {
				files = append(files, model.ConstructFile(file))
			}
			return nil
		})

		return files, err
	} else {
		f, err := os.Open(folder)
		if err != nil {
			return files, err
		}
		defer f.Close()

		if fileinfo, err := f.Readdir(-1); err == nil {
			pathSeparator := GetSeparator()
			for _, file := range fileinfo {
				if !file.IsDir() {
					folder, err := filepath.Abs(folder)
					if err != nil {
						return files, err
					}
					files = append(files, model.ConstructFile(folder+pathSeparator+file.Name()))
				} else {
					// fmt.Println("directory: ", file.Name())
				}
			}
		} else {
			return files, err
		}

	}

	return files, nil
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

// move file
func MoveFile(src, dest model.File) error {
	// Create the target folder at destination if not exists (for subfolder only)
	// Root target folder must exists, else error will be thrown. Validation is performed earlier on
	if _, err := os.Stat(dest.Folder); os.IsNotExist(err) {
		if err = os.MkdirAll(dest.Folder, os.ModePerm); err != nil {
			return err
		}
	}

	return os.Rename(src.FullPath, dest.FullPath)
}

// copy file
func CopyFile(src, dest model.File) error {
	source, err := os.Open(src.FullPath)
	if err != nil {
		return err
	}
	defer source.Close()

	// Create the target folder at destination if not exists (for subfolder only)
	// Root target folder must exists, else error will be thrown. Validation is performed earlier on
	if _, err := os.Stat(dest.Folder); os.IsNotExist(err) {
		if err = os.MkdirAll(dest.Folder, os.ModePerm); err != nil {
			return err
		}
	}

	destination, err := os.Create(dest.FullPath)

	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)

	return err
}
