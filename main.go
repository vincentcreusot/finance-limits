package main

import (
	"flag"
	"fmt"
	"github.com/vincentcreusot/finance-limits/fileutils"
	"github.com/vincentcreusot/finance-limits/logic"
	"os"
)

func main() {
	inputFileName := ""
	outputFileName := ""
	validateUsage(&inputFileName, &outputFileName)
	lineToParseChannel := make(chan string)
	go fileutils.ReadLines(inputFileName, lineToParseChannel)
	parser := logic.NewFinanceLogic()
	loadsToWrite, loadsErrors := parser.ParseLoads(lineToParseChannel)
	if len(loadsErrors) > 0 {
		for errCount, err := range loadsErrors {
			fmt.Printf("Error #%d in load: %v", errCount, err)
		}
	}
	err := fileutils.WriteLines(outputFileName, loadsToWrite)
	if err != nil {
		fmt.Println("Error writing lines:", err)
	}
}

func validateUsage(inputFileName *string, outputFileName *string) {
	flag.StringVar(inputFileName, "inputFile", "", "File to parse")
	flag.StringVar(inputFileName, "i", "", "File to parse")
	flag.StringVar(outputFileName, "outputFile", "", "File to write to")
	flag.StringVar(outputFileName, "o", "", "File to write to")
	flag.Parse()
	if *inputFileName == "" {
		fmt.Println("flag -inputFile is needed")
		flag.Usage()
		os.Exit(1)
	}
	if *outputFileName == "" {
		fmt.Println("flag -outputFile is needed")
		flag.Usage()
		os.Exit(1)
	}
}
