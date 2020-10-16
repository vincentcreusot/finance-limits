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
	outputFileName:= ""
	validateUsage(&inputFileName, &outputFileName)
	lineToParseChannel := make(chan string)
	go fileutils.ReadLines(inputFileName, lineToParseChannel)
	parser := logic.NewFinanceLogic()
	loadsToWrite,_ := parser.ParseLoads(lineToParseChannel)

	err := fileutils.WriteLines(outputFileName, loadsToWrite)
	if err != nil {
		fmt.Println ("Error ", err)
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
