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
	readUpperDir()

	getTerminalSize()

	maina()
}

func readCurrentDir() {
	readInDir(".")
}

func readUpperDir() {
	readInDir("..")
}


func readInDir(path string) {
	file, err := os.Open(path)
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

	calculateDirSize(path)
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
	fmt.Println("cwd: ", path)
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


func getTerminalSize()(cols int, rows int){
	cols, rows = consolesize.GetConsoleSize()
	fmt.Printf("Rows: %v, Cols: %v\n", rows, cols)
	return
}


func maina() {
    doneCh := make(chan struct{})

    signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH, syscall.SIGTERM)


    go receive(signalCh, doneCh)

	<-doneCh

	fmt.Println(doneCh)
}

func receive(signalCh chan os.Signal, doneCh chan struct{}) {
    for {
        select {
        // Example. Process to receive a message
        // case msg := <-receiveMessage():
		case sig := <-signalCh:
			doneCh <- struct{}{}
			// print this line telling us which signal was seen
			fmt.Println("Received signal from OS: ", sig)
        }
    }
}