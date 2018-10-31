package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbName := "goblog"
	dbPass := "T3$t!992"
	dbUser := "test"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	db := dbConn()
	defer db.Close()
}
