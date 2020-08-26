package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nathan-fiscaletti/consolesize-go"
)

func Open() {
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
