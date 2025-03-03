package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"

	"github.com/pi2pie/slugfiles/helper"
	"github.com/pi2pie/slugfiles/model"
)

var isRecursive bool
var isCaseSensitive bool

// Version can be set via ldflags during build
var Version = "0.0.4-beta.2"

// RootCmd is the root command for the CLI
var RootCmd = &cobra.Command{
	Use:     "slugfiles",
	Short:   "Rename files in a directory to user friendly slugs.",
	Version: Version,
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().StringP("output", "o", "", "Output directory")
	RootCmd.PersistentFlags().BoolVarP(&isRecursive, "recursive", "r", false, "Process directories recursively")
	RootCmd.PersistentFlags().BoolVarP(&isCaseSensitive, "case-sensitive", "c", false, "Case sensitive renaming")
	RootCmd.AddCommand(renameCmd)
}

// Rewritten outputFolder function to properly handle directory structures
func outputFolder() string {
	output, _ := RootCmd.PersistentFlags().GetString("output")
	if output == "" {
		return ""
	}
	
	// Normalize and ensure output has trailing separator
	output = filepath.Clean(output)
	if !strings.HasSuffix(output, string(os.PathSeparator)) {
		output += string(os.PathSeparator)
	}
	
	return output
}

// getTargetPath calculates the target path for a file in recursive mode
func getTargetPath(sourceBase, filePath, outputDir string) string {
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
	
	// Combine output directory with relative path
	return filepath.Join(outputDir, relPath) + string(os.PathSeparator)
}

var renameCmd = &cobra.Command{
	Use:   "rename [folder]",
	Short: "Rename files in a directory to user friendly slugs.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rename command called")

		// Check case sensitivity flag
		if isCaseSensitive {
			slug.Lowercase = false
		} else {
			slug.Lowercase = true
		}

		if args[0] != "" {
			// Clean and normalize the source directory path
			sourceDir := filepath.Clean(args[0])
			
			// Get files according to recursive flag
			files, err := helper.GetFiles(sourceDir, isRecursive)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("______________________")
			fmt.Println("Files in the folder: ", sourceDir)
			fmt.Println(" ")
			
			// Print files in tree structure using the helper function
			helper.PrintFilesTree(files, sourceDir)
			fmt.Println(" ")
			
			fmt.Println("______________________")
			fmt.Println("Renaming files...")
			fmt.Println(" ")
			
			// Get output directory
			outputDir := outputFolder()
			if outputDir != "" {
				fmt.Println("Output folder provided:", outputDir)
				// Check if the output folder exists, if not create it
				if !helper.HasDir(outputDir) {
					fmt.Println("Output folder does not exist")
					os.MkdirAll(outputDir, os.ModePerm)
					fmt.Println("Output folder created: ", outputDir)
				}
			}
			
			for _, file := range files {
				newname := slug.Make(file.FileName) + file.Ext
				
				// If original name is already slug-formatted, skip
				if newname == file.File {
					continue
				}
				
				// Handle outputDir differently if specified
				if outputDir != "" {
					var targetDir string
					if isRecursive {
						// Preserve directory structure under output folder
						targetDir = getTargetPath(sourceDir, file.FullPath, outputDir)
					} else {
						targetDir = outputDir
					}
					
					// Ensure target directory exists
					if !helper.HasDir(targetDir) {
						os.MkdirAll(targetDir, os.ModePerm)
					}
					
					newpath := filepath.Join(targetDir, newname)
					newfile := model.File{
						FullPath: newpath,
						Folder:   targetDir,
						File:     newname,
						FileName: strings.TrimSuffix(newname, file.Ext),
						Ext:      file.Ext,
					}
					
					// Copy the file with new name to output folder
					helper.CopyFile(file, newfile)
					fmt.Println(file.FullPath, "→", newpath)
				} else {
					// Rename the file in place
					newpath := filepath.Join(file.Folder, newname)
					newfile := model.File{
						FullPath: newpath,
						Folder:   file.Folder,
						File:     newname,
						FileName: strings.TrimSuffix(newname, file.Ext),
						Ext:      file.Ext,
					}
					
					helper.MoveFile(file, newfile)
					fmt.Println(file.FullPath, "→", newpath)
				}
			}
		}
	},
}

// Execute runs the root command
func Execute() error {
	return RootCmd.Execute()
}
