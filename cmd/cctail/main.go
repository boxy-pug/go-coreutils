package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type tailCommand struct {
	tailFiles []tailFile
	filePaths []string
	numLines  int
}

type tailFile struct {
	filePointer *os.File
	startPos    int64
}

func main() {
	var numLines int
	flag.IntVar(&numLines, "n", 10, "number of lines printed")
	flag.Parse()

	fileNames := flag.Args()

	tc := tailCommand{
		numLines:  numLines,
		filePaths: fileNames,
	}

	if len(tc.filePaths) == 0 {
		tc.tailFiles = append(tc.tailFiles, tailFile{filePointer: os.Stdin})
	} else {
		tc.openFile()
	}

	tc.getTail()

	tc.printTail()

}

func (tc *tailCommand) openFile() {
	for _, filepath := range tc.filePaths {
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("tail: %v\n", err)
			os.Exit(1)
		}
		tc.tailFiles = append(tc.tailFiles, tailFile{filePointer: file})
	}
}

func (tc *tailCommand) getTail() {
	for i := range tc.tailFiles {
		file := tc.tailFiles[i].filePointer

		pos, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			fmt.Println("Error seeking to the end of the file:", err)
			return
		}

		bufferSize := 4096
		buffer := make([]byte, bufferSize)
		lineCount := 0
		found := false
		var startPos int64

		if fileEndIsNewline(file, pos) {
			pos--
		}

		for pos > 0 && !found {
			if pos < int64(bufferSize) {
				bufferSize = int(pos)
			}
			pos -= int64(bufferSize)
			_, err := file.Seek(pos, io.SeekStart)
			if err != nil {
				fmt.Println("Error seeking in file:", err)
				return
			}
			//fmt.Printf("new pos is now: %v\n", pos)
			n, err := file.Read(buffer[:bufferSize])
			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}
			for j := n - 1; j >= 0; j-- {
				if buffer[j] == '\n' {
					//fmt.Printf("yey found newline at %v\n", pos+int64(j))
					lineCount++
					//fmt.Printf("Linecount is now: %v\n", lineCount)
					if lineCount == tc.numLines {
						//fmt.Printf("linecount is %v and tcnumlines is %v\n", lineCount, tc.numLines)
						// Found the starting point for the last N lines
						startPos = pos + int64(j+1)
						//fmt.Printf("startpos is: %v\n", startPos)
						found = true
						break
					}
				}
			}
		}
		if !found {
			startPos = 0
		}
		tc.tailFiles[i].startPos = startPos
	}
}

func (tc *tailCommand) printTail() {
	numberOfFiles := len(tc.tailFiles)
	for i, f := range tc.tailFiles {

		// Seek to the start position of the last N lines
		_, err := f.filePointer.Seek(f.startPos, io.SeekStart)
		if err != nil {
			fmt.Println("Error seeking in file:", err)
			return
		}

		// Read and print from startPos to the end of the file
		reader := bufio.NewReader(f.filePointer)
		data, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		if numberOfFiles > 1 {
			fmt.Printf("==> %s <==\n", f.filePointer.Name())
		}

		fmt.Printf("%s", string(data))
		if i < numberOfFiles-1 {
			fmt.Println()
		}
	}
}

func fileEndIsNewline(file *os.File, pos int64) bool {
	var lastByte []byte
	if pos > 0 {
		_, err := file.Seek(-1, io.SeekEnd)
		if err != nil {
			fmt.Println("Error seeking in file:", err)
			return false
		}
		lastByte = make([]byte, 1)
		_, err = file.Read(lastByte)
		if err != nil {
			fmt.Println("Error reading last byte:", err)
			return false
		}
	}
	return lastByte[0] == '\n'
}
