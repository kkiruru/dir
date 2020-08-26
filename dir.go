package main

import (

	"github.com/kkiruru/dir/model"
	"github.com/kkiruru/dir/view"

)


func main() {
	list := readCurrentDir()
	view.PrintDir(list)
}

func readCurrentDir() model.DirectoryInformation {
	return model.Read(".")
}

func changeDirectory(path string) model.DirectoryInformation {
	return model.Read(path)
}
