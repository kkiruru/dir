package view

import (
	"fmt"
	"os"

	"github.com/aybabtme/rgbterm"
	"github.com/kkiruru/dir/model"
)


func PrintDir(dirs model.DirectoryInformation) {

	var r, g, b uint8
	// pick a color
	r, g, b = 252, 255, 43
	// colorize it!
	coloredWord := rgbterm.FgString(dirs.Cwd, r, g, b)
	fmt.Println("> ", coloredWord)
	fmt.Println("Total ", dirs.Size)

	for _, dir := range dirs.SubDirs {
		fmt.Print(colorize(dir))
	}

	for _, file := range dirs.Files {
		fmt.Print(colorize(file))
	}

	fmt.Println("")
}

func colorize(file os.FileInfo) (coloredWord string) {
	var r, g, b uint8
	r, g, b = 252, 255, 43

	coloredWord = fmt.Sprintf("\n%s  %5s  %-23s ", model.GetDateTime(file), model.GetSize(file), file.Name())

	if file.IsDir() {
		if model.IsHidden(file) {
			r, g, b = 160, 0, 0
		} else {
			r, g, b = 224, 0, 0
		}
	} else if model.IsExecute(file) {
		if model.IsHidden(file) {
			r, g, b = 0, 160, 0
		} else {
			r, g, b = 0, 224, 0
		}
	} else {
		if model.IsHidden(file) {
			r, g, b = 160, 160, 160
		} else {
			r, g, b = 224, 224, 224
		}
	}

	coloredWord = rgbterm.FgString(coloredWord, r, g, b)
	return
}


func print(dirs []os.FileInfo) {
	fmt.Printf("\nName\t\tSize\tIsDirectory  Last Modification\n")
	for _, dirs := range dirs {
		fmt.Printf("\n%-15s  %-7v %-12v %v %#032b", dirs.Name(), dirs.Size(), dirs.IsDir(), dirs.ModTime(), dirs.Mode())
	}
}