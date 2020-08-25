package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/nathan-fiscaletti/consolesize-go"
)

type DirectoryInformation struct {
	cwd  string
	size int64
	// subDirs []FileInformation
	// files   []FileInformation
	subDirs []os.FileInfo
	files   []os.FileInfo
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
	result := readCurrentDir()
	result.cwd, _ = currentWorkingDirectory()
	printDir(result)
}

func readCurrentDir() DirectoryInformation {
	slice := []os.FileInfo{}
	fi, _ := dirInfo(".")
	slice = append(slice, fi)

	fi, _ = dirInfo("..")
	slice = append(slice, fi)

	dirs, files, size, _ := readInDir(".")
	slice = append(slice, dirs...)
	slice = append(slice, files...)

	var di DirectoryInformation

	di.size = size
	di.subDirs = dirs
	di.files = files

	return di
}

func printDir(dirs DirectoryInformation) {

	var r, g, b uint8
	// pick a color
	r, g, b = 252, 255, 43
	// colorize it!
	coloredWord := rgbterm.FgString(dirs.cwd, r, g, b)
	fmt.Println("> ", coloredWord)
	fmt.Println("Total ", dirs.size)

	for _, dir := range dirs.subDirs {
		fmt.Print(colorize(dir))
	}

	for _, file := range dirs.files {
		fmt.Print(colorize(file))
	}

	fmt.Println("")
}


func colorize(file os.FileInfo)(coloredWord string){
	var r, g, b uint8
	r, g, b = 252, 255, 43

	coloredWord = fmt.Sprintf("\n%-23s %-10s %s", file.Name(), getSize(file), getDateTime(file))

	if file.IsDir() {
		if isHidden(file){
			r, g, b = 160, 0, 0
		}else{
			r, g, b = 224, 0, 0
		}
	}else if isExecute(file) {
		if isHidden(file){
			r, g, b = 0, 160, 0
		}else{
			r, g, b = 0, 224, 0
		}
	}else{
		if isHidden(file){
			r, g, b = 160, 160, 160
		}else{
			r, g, b = 224, 224, 224
		}
	}

	coloredWord = rgbterm.FgString(coloredWord, r, g, b)

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

func isExecute(file os.FileInfo) bool {
	return ( ((uint32(file.Mode()) & (OTHER_X|GROUP_X|OWNER_X)) != 0 ))
}

func isHidden(file os.FileInfo) bool {
	name := file.Name()
	return name[0] == '.'
}

func getSize(file os.FileInfo) string {

	fileSize := file.Size()

	if fileSize < 1024 {
		return strconv.FormatUint(uint64(fileSize), 10)
	} else {
		fileSize = fileSize / 1024
		if fileSize < 1024 {
			return strconv.FormatUint(uint64(fileSize), 10)+"K"
		} else {
			fileSize = fileSize / 1024
			if fileSize < 1024 {
				return strconv.FormatUint(uint64(fileSize), 10)+"M"
			} else {
				fileSize = fileSize / 1024
				if fileSize < 1024 {
					return strconv.FormatUint(uint64(fileSize), 10)+"G"
				}
			}
		}
		return string(fileSize)
	}
}

func getDateTime(file os.FileInfo) string {
	t := file.ModTime()
	return t.Format("2006-01-02 15:04:05")
}




func print(dirs []os.FileInfo) {
	fmt.Printf("\nName\t\tSize\tIsDirectory  Last Modification\n")
	for _, dirs := range dirs {
		fmt.Printf("\n%-15s %-7v %-12v %v %#032b", dirs.Name(), dirs.Size(), dirs.IsDir(), dirs.ModTime(), dirs.Mode())
	}
}

func readUpperDir() {
	readInDir("..")
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
		return files[i].Name() < files[j].Name()
	})

	return
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
	return info, err
}
