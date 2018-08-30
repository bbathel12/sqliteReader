package main

import (
	"database/sql"
	"fmt"
	"os"
)

func readRows(db *sql.DB, hashChan *chan string) {
	rows, err := db.Query("select * from emails ")
	if err != nil {
		handleErr(err)
	}

	for rows.Next() {
		var email string
		rows.Scan(&email)
		*hashChan <- email
	}
	close(*hashChan)
}

func writeToFile(outputDir string, hashChan *chan string) { //, wg *sync.WaitGroup) {
	file, err := os.Create(outputDir + "/output.txt")
	handleErr(err)
	defer file.Close()

	for str := range *hashChan {
		bytes, err := file.WriteString(str + "\n")
		handleErr(err)
		if err != nil {
			fmt.Println(bytes)
		}
		//writer.Flush()
	}
}
