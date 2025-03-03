package main

import (
	"fmt"
	"os"
)

var application applicationType

func main() {
	if len(os.Args) != 2 {
		fmt.Println("No log-file name sent")
		os.Exit(1)
	} else {
		application.init()
	}
}
