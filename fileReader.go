package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

func readUploadByLine(uploadName string, lineChan *chan string, readerGroup *sync.WaitGroup) {
	upload, err := os.Open(uploadName)
	defer upload.Close()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(upload)
	scanner.Buffer([]byte{}, 200)

	for scanner.Scan() {
		*lineChan <- scanner.Text()
	}

	defer readerGroup.Done()

}

func readUploadDirectory(uploadDirectory string, lineChan *chan string) {
	var readerGroup sync.WaitGroup

	files, err := ioutil.ReadDir(uploadDirectory)

	if err != nil {
		fmt.Println("Input Directory", uploadDirectory, "Not found")
		os.Exit(1)
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			readerGroup.Add(1)
			go readUploadByLine(uploadDirectory+"/"+file.Name(), lineChan, &readerGroup)
		}

	}

	readerGroup.Wait()

	defer close(*lineChan)
}
