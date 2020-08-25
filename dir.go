package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nathan-fiscaletti/consolesize-go"
)

type DirectoryInformation struct {
	cwd     string
	size    int64
	subDirs []FileInformation
	files   []FileInformation
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
	// directoryName, _ := currentWorkingDirectory()
	// readInDir(directoryName)

	directoryName, _ := currentWorkingDirectory()
	fmt.Println("> ", directoryName)

	// dirInfo(".")
	// dirInfo("..")

	readCurrentDir()
	// readUpperDir()
	// getTerminalSize()
}

func readCurrentDir() {
	slice := []os.FileInfo{}
	fi, _ := dirInfo(".")
	slice = append(slice, fi)

	fi, _ = dirInfo("..")
	slice = append(slice, fi)

	dirs, _ := readInDir(".")
	slice = append(slice, dirs...)


	dirs = []os.FileInfo{}
	for _, dir := range slice {

		if dir.IsDir() {
			dirs = append(dirs, dir)
		}
	}

	for _, dir := range slice {
		if !dir.IsDir() {
			dirs = append(dirs, dir)
		}
	}
	print(dirs)
}

func print(dirs []os.FileInfo){

	fmt.Printf("\nName\t\tSize\tIsDirectory  Last Modification\n")
	for _, dirs := range dirs {
		fmt.Printf("\n%-15s %-7v %-12v %v %#032b", dirs.Name(), dirs.Size(), dirs.IsDir(), dirs.ModTime(), dirs.Mode())
	}
}


func readUpperDir() {
	readInDir("..")
}

func readInDir(path string) ([]os.FileInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	fileList, _ := file.Readdir(0)

	// fmt.Printf("\nName\t\tSize\tIsDirectory  Last Modification\n")
	// for _, files := range fileList {
	// 	fmt.Printf("\n%-15s %-7v %-12v %v %#032b", files.Name(), files.Size(), files.IsDir(), files.ModTime(), files.Mode())
	// }

	// fmt.Println()

	return fileList, nil
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
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	// fmt.Println("cwd: ", path)
	return path, nil
}

func calculateDirSize(path string) (size int64, err error) {
	isDir, _ := isDirectory(path)
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

func isDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), err
}

func getTerminalSize() (cols int, rows int) {
	cols, rows = consolesize.GetConsoleSize()
	fmt.Printf("Rows: %v, Cols: %v\n", rows, cols)
	return
}

func dectectionTerminalSize() {
	doneCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, syscall.SIGWINCH, syscall.SIGTERM)

	go receive(signalCh, doneCh)

	<-doneCh

	fmt.Println("doneCh")

	getTerminalSize()
}

func receive(signalCh chan os.Signal, doneCh chan struct{}) {
	for {
		select {
		case sig := <-signalCh:
			fmt.Println("Received signal from OS: ", sig)
			doneCh <- struct{}{}
		}
	}
}

func dirInfo(dir string) (os.FileInfo, error) {
	file, err := os.Open(dir)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	info, _ := file.Stat()

	// fmt.Printf("%-15s %-7v %-12v %v %#032b\n", info.Name(), info.Size(), info.IsDir(), info.ModTime(), info.Mode())

	return info, err
}
