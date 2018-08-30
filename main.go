//Package to write the output file because php keeps crashing on me
package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type listArgs struct {
	UploadDir string
	S3Dir     string
	DBDir     string
	OutputDir string
}

// this is will be the buffer size for all chans
var buffersize int = 10000

// make all channels with buffersize
var lineChan chan string = make(chan string, buffersize)
var domainChan chan string = make(chan string, buffersize)
var hashedLineChan chan string = make(chan string, buffersize)
var errorChan chan string = make(chan string, buffersize)
var newHashChan chan string = make(chan string, buffersize)

//unbuffered channel just for blocking;
var doneChan chan bool = make(chan bool)

func getArgs() (dirs listArgs) {

	args := os.Args[1:]

	if len(args) < 4 {
		fmt.Println("Usage: sqlReader <directory>")
		os.Exit(1)
	}

	dirs.UploadDir = args[0]
	dirs.S3Dir = args[1]
	dirs.DBDir = args[2]
	dirs.OutputDir = args[3]

	return
}

func main() {
	var dirs listArgs = getArgs()
	var recs int = 0
	start := time.Now()
	// create database connection
	db, err := openConnection(dirs.DBDir + "/dupes.db")
	handleErr(err)

	// create tables for deduping
	createTables(db)

	//read in the entire upload directory and send out each line on lineChan
	go readUploadDirectory(dirs.UploadDir, &lineChan)

	// multiplex line chan into domain,md5 and error chans
	go domainEmailMultiPlex(&recs, &errorChan, &lineChan, &domainChan, &hashedLineChan)

	// insert new hashes that are in hashedLineChan
	go insertHashes(db, &hashedLineChan, &doneChan)

	<-doneChan

	go readRows(db, &newHashChan)

	//Open file here to pass in
	//this is incase I decide to go back to wait groups
	//file, err := os.Create(dirs.OutputDir + "/output.txt")

	writeToFile(dirs.OutputDir, &newHashChan) //, &wg)
	fmt.Println("DEDUPING: ", time.Since(start))
	//file.Close()

}
