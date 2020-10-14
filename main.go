package main

import (
	"flag"
	"fmt"
	"github.com/vincentcreusot/finance-limits/fileutils"
	"os"
)

func main() {
	inputFileName := ""
	validateUsage(&inputFileName)
	fmt.Println("File ", inputFileName)
	lineToParseChannel := make(chan string)
	go fileutils.ReadLines(inputFileName, lineToParseChannel)
	for line := range lineToParseChannel {
		fmt.Println(line)
	}

}

func validateUsage(inputFileName *string) {
	flag.StringVar(inputFileName, "inputFile", "", "File to parse")
	flag.StringVar(inputFileName, "i", "", "File to parse")
	flag.Parse()
	fmt.Println("File ", *inputFileName)
	if *inputFileName == "" {
		fmt.Println("flag -inputFile is needed")
		flag.Usage()
		os.Exit(1)
	}
}
