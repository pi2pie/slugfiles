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
var isDryRun bool

// Version can be set via ldflags during build
var Version = "0.0.4-rc.5"

// RootCmd is the root command for the CLI
var RootCmd = &cobra.Command{
    Use:     "slugfiles [folder]",
    Short:   "Rename files in a directory to user friendly slugs.",
    Version: Version,
    Args:    cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) == 1 {
            // If a file path is provided, execute the rename command's logic
            renameArgs := []string{args[0]}
            renameCmd.Run(cmd, renameArgs)
        } else {
            // No argument: display help
            cmd.Help()
        }
    },
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().StringP("output", "o", "", "Output directory")
	RootCmd.PersistentFlags().BoolVarP(&isRecursive, "recursive", "r", false, "Process directories recursively")
	RootCmd.PersistentFlags().BoolVarP(&isCaseSensitive, "case-sensitive", "c", false, "Case sensitive renaming")
	RootCmd.PersistentFlags().BoolVarP(&isDryRun, "dry-run", "d", false, "Simulate renaming without making changes")
	RootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:   "rename [folder]",
	Short: "Rename files in a directory to user friendly slugs.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rename command called")

		// dealing with the output folder
		outputFolder := func() string {
			outputFlag, _ := cmd.Flags().GetString("output")
			if outputFlag == "" {
				return ""
			}
			// Normalize and ensure output has trailing separator
			output := filepath.Clean(outputFlag)
			if !strings.HasSuffix(output, string(os.PathSeparator)) {
				output += string(os.PathSeparator)
			}
			
			return output
		}

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
					if !isDryRun {
						os.MkdirAll(outputDir, os.ModePerm)
						fmt.Println("Output folder created: ", outputDir)
					} else {
						fmt.Println("[DRY RUN] Would create output folder: ", outputDir)
					}
				}
			}
			
			// Process files
			for _, file := range files {
				newname := slug.Make(file.FileName) + file.Ext
				
				// Check if file needs renaming
				needsRenaming := newname != file.File
    
				// Only skip if the file doesn't need renaming AND we're not copying to output folder
				if !needsRenaming && outputDir == "" {
					continue
				}
				
				// Handle outputDir differently if specified
				if outputDir != "" {
					var targetDir string
					if isRecursive {
						// Preserve directory structure under output folder
						targetDir = helper.GetSlugifiedTargetPath(sourceDir, file.Folder, outputDir, slug.Make)
					} else {
						targetDir = outputDir
					}
					
					// Inside the file processing loop where directories are created:
                    // Ensure target directory exists
                    if !helper.HasDir(targetDir) {
                        if !isDryRun {
                            os.MkdirAll(targetDir, os.ModePerm)
                        }
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
                    if !isDryRun {
                        helper.CopyFile(file, newfile)
                    }
        
					if needsRenaming {
						if isDryRun {
                            fmt.Println("[DRY RUN] Would rename:", file.FullPath, "→", newpath)
                        } else {
                            fmt.Println(file.FullPath, "→", newpath)
                        }
					} else {
						if isDryRun {
                            fmt.Println("[DRY RUN] Would copy:", file.FullPath, "→", newpath)
                        } else {
                            fmt.Println(file.FullPath, "→ (copied to)", newpath)
                        }
					}
				} else {
					// Only rename if needed (name is not already slug-formatted)
					if needsRenaming {
						// Rename the file in place
						newpath := filepath.Join(file.Folder, newname)
						newfile := model.File{
							FullPath: newpath,
							Folder:   file.Folder,
							File:     newname,
							FileName: strings.TrimSuffix(newname, file.Ext),
							Ext:      file.Ext,
						}
						
						if !isDryRun {
                            helper.MoveFile(file, newfile)
                            fmt.Println(file.FullPath, "→", newpath)
                        } else {
                            fmt.Println("[DRY RUN] Would rename:", file.FullPath, "→", newpath)
                        }
					}
				}
			}
			// If recursive and no output directory, rename directories
            if isRecursive && outputDir == "" {
                fmt.Println("______________________")
                if isDryRun {
                    fmt.Println("[DRY RUN] Preview of directory renaming (no changes will be made)...")
                } else {
                    fmt.Println("Renaming directories...")
                }
                fmt.Println(" ")
                
                // Get unique directories
                dirs := helper.GetDirectories(files, sourceDir)
                
                // Sort directories by depth (deepest first)
                helper.SortDirsByDepth(dirs)
                
                // Rename directories
                for _, dir := range dirs {
                    dirName := filepath.Base(dir)
                    parentDir := filepath.Dir(dir)
                    
                    // Create slug for directory name
                    slugDirName := slug.Make(dirName)
                    
                    if slugDirName != dirName {
                        // we make a new directory with the slug name
						newDir := filepath.Join(filepath.Dir(parentDir), slugDirName)
						if isDryRun {
							fmt.Println("[DRY RUN] Would rename directory:", dir, "→", newDir)
						} else {
							os.MkdirAll(newDir, os.ModePerm)
							fmt.Println(dir, "→", newDir)
							// then we move all files from the old directory to the new one
							if err := helper.MoveFilesByPath(dir, newDir); err != nil {
								fmt.Println("Error moving files:", err)
							}
							// and finally we remove the old directory
							if err := os.RemoveAll(dir); err != nil {
								fmt.Println("Error removing old directory:", err)
							} else {
								fmt.Println("Removed old directory:", dir)
							}
						}
                    }
                }
            }

		}
	},
}

// Add this special handling to Execute function
func Execute() error {
    // Special case handling for "rename" directory
    if len(os.Args) >= 2 && os.Args[1] == "rename" {
        // Check if "rename" is a directory
        info, err := os.Stat("rename")
        if err == nil && info.IsDir() {
            // It's a directory, so treat it as a path argument rather than a subcommand
            
            // Capture any flags that were provided after "rename"
            var flags []string
            for i := 2; i < len(os.Args); i++ {
                if strings.HasPrefix(os.Args[i], "-") {
                    flags = append(flags, os.Args[i])
                }
            }
            
            // Parse flags first
            if len(flags) > 0 {
                RootCmd.ParseFlags(os.Args[2:])
            }
            
            // Pass the directory as argument
            renameArgs := []string{"rename"}
            renameCmd.Run(RootCmd, renameArgs)
            return nil
        }
    }
    
    return RootCmd.Execute()
}