package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type DirectoryInfo struct {
	cwd   string
	size  int64
	files []FileInformation
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

func main() {
	fmt.Println("called main")

	// directoryName, _ := currentWorkingDirectory()
	// readInDir(directoryName)
	readCurrentDir()
	upperDirInfo()
}

func readCurrentDir() {
	readInDir(".")
}

func upperDirInfo() {
	readInDir("..")
}


func readInDir(directory string) {
	file, err := os.Open(directory)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	fileList, _ := file.Readdir(0)

	fmt.Printf("\nName\t\tSize\tIsDirectory  Last Modification\n")

	for _, files := range fileList {
		fmt.Printf("\n%-15s %-7v %-12v %v %#b", files.Name(), files.Size(), files.IsDir(), files.ModTime(), files.Mode())
	}

	fmt.Println()

	calculateDirSize(directory)
}

func waitingKeyEvent() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		for {
			timer := time.NewTimer(time.Second * 1)
			<-timer.C
			fmt.Println(scanner.Text())
		}
	}
	if scanner.Err() != nil {
		/*handle error*/
	}
}

func currentWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println("cwd: ", dir)
	return dir, nil
}

func calculateDirSize(dirpath string) (dirsize int64, err error) {
	isDir, _ := isDirectory(dirpath)
	if !isDir {
		return
	}

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.Mode().IsRegular() {
			dirsize += file.Size()
		}
	}

	fmt.Printf("_ %s size is %d \n", dirpath, dirsize)
	return
}

func isDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), err
}
