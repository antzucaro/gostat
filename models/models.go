package models

import (
	"database/sql"
	"github.com/antzucaro/gostat/config"
	_ "github.com/lib/pq"
	"log"
)

// the main connection used all throughout the app
var db *sql.DB

// Init opens a connection to the database and prepares all of the queries
func Init() (err error) {
	// establish a database connection
	db, err = sql.Open("postgres", config.Config.ConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// connection pooling
	db.SetMaxIdleConns(5)

	// prepare all of the queries

	// leaderboard
	topNRanksStmt = initStatement("topNRanksStmt", topNRanksSQL)
	recentActivePlayersStmt = initStatement("recentActivePlayersStmt", recentActivePlayersSQL)
	overallActivePlayersStmt = initStatement("overallActivePlayersStmt", overallActivePlayersSQL)
	recentGameCountStmt = initStatement("recentGameCountStmt", recentGameCountSQL)
	overallGameCountStmt = initStatement("overallGameCountStmt", overallGameCountSQL)
	playerTimeStmt = initStatement("playerTimeStmt", playerTimeSQL)
	serverGamesStmt = initStatement("serverGamesStmt ", serverGamesSQL)
	mapGamesStmt = initStatement("mapGamesStmt ", mapGamesSQL)
	playerScoreStmt = initStatement("playerScoreStmt", playerScoreSQL)
	recentGamesStmt = initStatement("recentGamesStmt", recentGamesSQL)

	return
}

// initializes a prepared statement by name for better logging/traceability
func initStatement(name string, sql string) *sql.Stmt {
	stmt, err := db.Prepare(sql)
	if err != nil {
		//log.Fatal("Error preparing SQL statement " + name)
		log.Fatal(err)
	} else {
		log.Println("Prepared statement " + name)
	}

	return stmt
}
