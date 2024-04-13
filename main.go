package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	folderstable "github.com/erancihan/remove-dep-folders/internal/folders-table"
	"github.com/erancihan/remove-dep-folders/internal/utils"
	"github.com/spf13/cobra"
)

var (
	isDryRun bool = false
	checkAll bool = false
)

func getFolders(searchpath string) ([]folderstable.FolderEntry, error) {
	entries := []folderstable.FolderEntry{}

	err := filepath.WalkDir(searchpath, func(path string, d fs.DirEntry, err error) error {
		// check if directory
		if !d.IsDir() {
			return nil
		}

		// ignore dot directories, except for the current directory
		if !checkAll && d.Name() != "." && len(d.Name()) > 1 && d.Name()[0] == '.' {
			return filepath.SkipDir
		}

		// check if node_modules directory
		if d.Name() == "node_modules" {
			size, err := utils.DirSizeB(path)
			if err != nil {
				fmt.Println(err)
			}

			entries = append(entries, folderstable.FolderEntry{Path: path, Size: size})

			// do not recurse into node_modules
			return filepath.SkipDir
		}

		// check if python venv directory
		// folder should contain pyvenv.cfg file to be considered a venv
		if _, err := os.Stat(filepath.Join(path, "pyvenv.cfg")); err == nil {
			size, err := utils.DirSizeB(path)
			if err != nil {
				fmt.Println(err)
			}

			entries = append(entries, folderstable.FolderEntry{Path: path, Size: size})

			// do not recurse into venv
			return filepath.SkipDir
		}

		return nil
	})

	return entries, err
}

func removeFolders(paths []string) {
	for _, path := range paths {
		fmt.Printf("Removing %s\n", path)

		// sleep for a second to allow the user to cancel the operation
		if isDryRun {
			continue
		}

		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func process(searchpath string) {
	entries, err := getFolders(searchpath)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(entries) == 0 {
		fmt.Println("No folders found")
		return
	}

	selection := folderstable.ShowFolderSelection(entries)
	if len(selection) == 0 {
		fmt.Println("No folders selected")
		return
	}

	removeFolders(selection)
}

var cmd = &cobra.Command{
	Use:   "remove-dep-folders <path>",
	Short: "Remove dependency folders from your system",
	Long:  "A tool to remove node_modules and vendor directories from your system",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string

		if len(args) > 0 {
			path = args[0]
		}
		if path == "" {
			cmd.Help()
			return
		}

		// get the absolute path
		path, err := filepath.Abs(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if the path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Path '%s' does not exist\n\n", path)
			cmd.Help()
			return
		}

		process(path)
	},
}

func init() {
	cmd.PersistentFlags().BoolVar(&isDryRun, "dry-run", false, "Dry run the program")
	cmd.PersistentFlags().BoolVarP(&checkAll, "all", "a", false, "Check all folders")
}

func main() {
	cmd.Execute()
}
