package fileutils

import (
	"bufio"
	"log"
	"os"
)

// ReadLines read a file and send each line to a channel
func ReadLines(inputFileName string, lineChannel chan string) {
	defer close(lineChannel)
	fileBuffer, err := os.Open(inputFileName)
	if err != nil {
		log.Fatal("Error opening file ", inputFileName, err)
	}

	defer func() {
		if err = fileBuffer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	lineScanner := bufio.NewScanner(fileBuffer)
	for lineScanner.Scan() {
		lineChannel <- lineScanner.Text()
	}
	err = lineScanner.Err()
	if err != nil {
		log.Fatal(err)
	}
}

