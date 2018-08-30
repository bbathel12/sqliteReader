package main

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func handleErr(err error) {
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
}

func openConnection(dbFilePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:"+dbFilePath+"?cache=shared&mode=memory&_jounal=off&_sync=0") //dbFilePath)
	handleErr(err)
	return db, err
}

func createTables(db *sql.DB) {
	_, err := db.Exec("create table emails (email text not null)")
	handleErr(err)

	_, err = db.Exec("create unique index unique_email on emails(email)")
	handleErr(err)

	_, err = db.Exec("PRAGMA synchronous = OFF;")
	handleErr(err)

	_, err = db.Exec("PRAGMA jounal_mode = OFF;")
	handleErr(err)

	_, err = db.Exec("PRAGMA cache_size=10000;")
	handleErr(err)

}

func insertHashes(db *sql.DB, hashChan *chan string, doneChan *chan bool) {
	// number of items to batch insert at a time for fine tuning
	const batchAmount int = 999

	// track number for batch insert
	var values int = 0

	// this will store the strings for batch insert
	//args := [batchAmount]string{}
	var args []interface{}

	stmt, err := db.Prepare(buildBulkStatement(batchAmount))
	handleErr(err)

	tx, err := db.Begin()
	handleErr(err)

	for str := range *hashChan {
		args = append(args, str)
		values++

		if values%batchAmount == 0 {

			// run statement with args handle errors
			_, err = stmt.Exec(args...)
			handleErr(err)

			// reset everything
			values = 0
			args = args[:0]
		}
	}
	//tx.Commit()

	//tx, err = db.Begin()
	//handleErr(err)

	stmt, err = db.Prepare(buildBulkStatement(values))
	handleErr(err)

	_, err = stmt.Exec(args...)
	handleErr(err)

	tx.Commit()

	// Close everything down
	close(*doneChan)
}

func buildBulkStatement(NumberOfParams int) string {
	var questionMarks []string

	for i := 0; i < NumberOfParams; i++ {
		questionMarks = append(questionMarks, "(?)")
	}

	// using strings.Builder to create prepared statement for batch
	var sql strings.Builder
	var sqlStart string = "insert or ignore into emails (email) values "
	sql.WriteString(sqlStart)
	sql.WriteString(strings.Join(questionMarks, ","))
	return sql.String()
}
