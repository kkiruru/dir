package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)


type DirectoryInformation struct {
	Cwd string
	Size int64
	// subDirs []FileInformation
	// files   []FileInformation
	SubDirs []os.FileInfo
	Files   []os.FileInfo
}

type FileInformation struct {
	name      string
	size      int64
	date      string
	time      string
	isHidden  bool
	isExecute bool
	isDir     bool
}


func Read(path string) DirectoryInformation {

	slice := []os.FileInfo{}

	dot, _ := dirInfo(".")
	slice = append(slice, dot)

	dotdot, _ := dirInfo("..")
	slice = append(slice, dotdot)

	var di DirectoryInformation
	di.SubDirs, di.Files, di.Size, _ = readInDir(path)

	return di
}

func presentWorkingDirectory() (string) {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return "."
	}
	return path
}

func readInDir(path string) (dirs []os.FileInfo, files []os.FileInfo, size int64, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	fileList, _ := file.Readdir(0)

	dirs = []os.FileInfo{}
	files = []os.FileInfo{}
	size = 0

	current, _ := dirInfo(".")
	dirs = append(dirs, current)

	parent, _ := dirInfo("..")
	dirs = append(dirs, parent)

	for _, file := range fileList {
		if file.IsDir() {
			dirs = append(dirs, file)
		} else {
			files = append(files, file)
		}
		size = size + file.Size()
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	sort.Slice(files, func(i, j int) bool {
		return strings.ToUpper(files[i].Name()) < strings.ToUpper(files[j].Name())
	})

	return
}

const (
	OTHER_X uint32 = 1 << iota
	OTHER_W
	OTHER_R
	GROUP_X
	GROUP_W
	GROUP_R
	OWNER_X
	OWNER_W
	OWNER_R
)

func IsExecute(file os.FileInfo) bool {
	return ((uint32(file.Mode()) & (OTHER_X | GROUP_X | OWNER_X)) != 0)
}

func IsHidden(file os.FileInfo) bool {
	name := file.Name()
	return name[0] == '.'
}

func GetSize(file os.FileInfo) string {

	fileSize := file.Size()

	if fileSize < 1024 {
		return strconv.FormatUint(uint64(fileSize), 10)
	} else {
		fileSize = fileSize / 1024

		if fileSize < 1024 {
			return strconv.FormatUint(uint64(fileSize), 10) + "K"
		} else {
			fileSize = fileSize / 1024
			if fileSize < 1024 {
				return strconv.FormatUint(uint64(fileSize), 10) + "M"
			} else {
				fileSize = fileSize / 1024
				if fileSize < 1024 {
					return strconv.FormatUint(uint64(fileSize), 10) + "G"
				}
			}
		}
		return fmt.Sprint(fileSize)
	}
}

func GetDateTime(file os.FileInfo) string {
	t := file.ModTime()
	return t.Format("2006-01-02 15:04:05")
}

func calculateDirSize(path string) (size int64, err error) {
	isDir, _ := IsDirectory(path)
	if !isDir {
		return
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			size += file.Size()
		}
	}

	fmt.Printf("_ %s size is %d \n", path, size)
	return
}

func IsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), err
}

func dirInfo(dir string) (os.FileInfo, error) {
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	info, _ := file.Stat()
	return info, err
}
