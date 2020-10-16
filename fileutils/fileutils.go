package fileutils

import (
	"bufio"
	"fmt"
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

func WriteLines(filename string, loadsToWrite []string) error {
	if fileExists(filename) {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	for _,line := range loadsToWrite {
		fmt.Fprintln(f, line)
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}


func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}