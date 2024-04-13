package folderstable

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/erancihan/remove-dep-folders/internal/utils"
)

type FolderEntry struct {
	Path string
	Size int64
}

func ShowFolderSelection(choices []FolderEntry) []string {
	maxSizeWidth := 0
	for _, choice := range choices {
		// find the size text width
		sizeText := utils.ByteCounSI(choice.Size)
		if len(sizeText) > maxSizeWidth {
			maxSizeWidth = len(sizeText)
		}
	}

	// sort choices by size
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].Size > choices[j].Size
	})

	items := make([]huh.Option[string], len(choices))
	for i, choice := range choices {
		items[i] = huh.NewOption(
			fmt.Sprintf("%*s %s", maxSizeWidth, utils.ByteCounSI(choice.Size), choice.Path),
			choice.Path,
		)
	}

	selection := []string{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Options(items...).
				Title(("Select folders to delete")).
				Value(&selection).
				Height(15),
		),
	)

	if err := form.Run(); err != nil {
		if err == huh.ErrUserAborted {
			fmt.Println("Program was interrupted by the user")
			os.Exit(0)
		}

		fmt.Printf("Error running program: %v ", err)
		os.Exit(1)
	}

	return selection
}
