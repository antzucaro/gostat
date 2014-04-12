package db

import (
    _ "github.com/lib/pq"
    "database/sql"
    "log"
)

var db *sql.DB

func Init(connString string) (err error) {
  // establish a database connection
  db, err = sql.Open("postgres", connString)
  if err != nil {
    log.Fatal(err)
  }

  // prepare all of the queries
  topNRanksStmt, err = db.Prepare(topNRanksSQL)

  return
}
