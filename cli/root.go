package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"

	"github.com/pi2pie/slugfiles/helper"
	"github.com/pi2pie/slugfiles/model"
)

var isRecursive bool

// RootCmd is the root command for the CLI
var RootCmd = &cobra.Command{
	Use:   "slugfiles",
	Short: "Rename files in a directory to user friendly slugs.",
	Version: "0.0.2",	
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().BoolVarP(&isRecursive, "recursive", "r", false, "Recursively rename files in subdirectories")
	RootCmd.PersistentFlags().StringP("output", "o", "", "Output directory")
	RootCmd.AddCommand(renameCmd)
}

func outputFolder(path string) string {
	output, _ := RootCmd.PersistentFlags().GetString("output")
	// seperate path to the array seperated by /
	pathArray := strings.Split(path, "/")
	// remove the last element
	// if pathArray is empty skip the operation
	if len(pathArray) == 0 {
		return output
	} else if len(pathArray) == 1 {
		return output
	} else {
	pathArray = pathArray[:len(pathArray)-2]
	// fmt.Println("afeter: ", pathArray)
	// join the array back to a string
	output = strings.Join(pathArray, "/") + "/" + output
	return output
	}
}

var renameCmd = &cobra.Command{
	Use:   "rename [folder]",
	Short: "Rename files in a directory to user friendly slugs.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rename command called")

		if args[0] != "" {
			// if args[0] is a folder, print the files in the folder
			// read the folder and loop through the files and print the file names
			files, err := helper.GetFiles(args[0], isRecursive)			
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("______________________")
			fmt.Println(" ")
			for _, file := range files {
				fmt.Println(file.FullPath)
							
			}
			
			fmt.Println("______________________")
			fmt.Println(" ")
			if outputFolder("") != "" {
				fmt.Println("Output folder provided")
				// check the output folder exists, if not create the folder
				if !helper.HasDir(outputFolder(args[0])) {
					fmt.Println("Output folder does not exist")
					os.MkdirAll(outputFolder(args[0]), os.ModePerm)
					fmt.Println("Output folder created: ", outputFolder(args[0]))
				}
			}
			for _, file := range files {
				newname := slug.Make(file.FileName) + file.Ext
				newpath := file.Folder + newname				
				newfile := model.File {
					FullPath: newpath,
					Folder: file.Folder,
					File: newname,
					FileName: file.FileName,
					Ext: file.Ext,
				}				
				// if output folder is provided, move the file to the output folder
				if outputFolder("") != "" {
					// fmt.Println(outputFolder(file.Folder))					
					newpath = outputFolder(file.Folder) + "/" + newname
					newfile = model.File {
						FullPath: newpath,
						Folder: outputFolder(file.Folder),
						File: newname,
						FileName: file.FileName,
						Ext: file.Ext,
					}
					// replace the file name with the new name
				if newname != file.File {
					// exec.Command("cp", file.FullPath, newpath).Run()
					helper.CopyFile(file, newfile)
					fmt.Println(newpath)
					}					
				} else {
					// replace the file name with the new name
					if newname != file.File {
						exec.Command("mv", file.FullPath, newpath).Run()
						// helper.MoveFile(file, newfile)
						fmt.Println(newpath)
					}
				}				
			}
		}
	},
}

// Execute runs the root command
func Execute() error {
	return RootCmd.Execute()
}
