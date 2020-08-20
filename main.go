package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("called main")

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
