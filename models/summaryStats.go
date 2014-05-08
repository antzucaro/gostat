package models

import (
	"database/sql"
	"github.com/antzucaro/gostat/config"
	"log"
)

// SummaryStats represents summary information for a period of time:
//
// - Number of active players
// - Number of games played
// - A breakdown of the number of games played by the top five game types  The
//   other game types played but not in the top five are grouped into "other".
type SummaryStats struct {
	Players    int
	Games      int
	GameCounts []GameCount
	OtherGames int
}

var recentActivePlayersSQL = `SELECT count(distinct player_id) 
FROM player_game_stats 
WHERE player_id > 1
AND create_dt >= now() at time zone 'utc' - interval '` +
	config.Config.SummaryStatsDays + " days'"

var recentActivePlayersStmt *sql.Stmt

const overallActivePlayersSQL = `SELECT count(distinct player_id) 
FROM players 
WHERE player_id >= 2 and active_ind = 'Y'`

var overallActivePlayersStmt *sql.Stmt

// Retrieves the number of active players, with a flag "recent" to pull
// all or just the most recent. Bots are excluded, but one anonymous player
// is included - we cannot count the distinct number of those types of
// players (they are all player_id 2).
func GetActivePlayers(recent bool) int {
	var row *sql.Row
	if recent {
		row = recentActivePlayersStmt.QueryRow()
	} else {
		row = overallActivePlayersStmt.QueryRow()
	}

	var players int
	err := row.Scan(&players)
	if err != nil {
		log.Fatal(err)
	}

	return players
}

// GameCount is a struct to hold the various game types and the
// number of times they have been played within a given interval
type GameCount struct {
	GameType string
	Games    int
}

const overallGameCountSQL = `SELECT game_type_cd, count(*) 
FROM games 
GROUP BY game_type_cd
ORDER BY count(*) desc`

var overallGameCountStmt *sql.Stmt

var recentGameCountSQL = `SELECT game_type_cd, count(*) 
FROM games 
WHERE create_dt >= now() at time zone 'utc' - interval '` +
	config.Config.SummaryStatsDays + " days' " +
	`GROUP BY game_type_cd
ORDER BY count(*) desc`

var recentGameCountStmt *sql.Stmt

// Retrieve the games played in the past 24 hours by game type.
// Returns an array of GameCount structs ordered from the most
// played game type to the least. No aggregations are performed.
func GetGameCounts(recent bool) []GameCount {
	var rows *sql.Rows
	var err error
	if recent {
		// fetch the count of games from "days" number of days ago and later
		rows, err = recentGameCountStmt.Query()
	} else {
		// fetch the total count of games
		rows, err = overallGameCountStmt.Query()
	}

	if err != nil {
		log.Fatal(err)
	}

	gameCounts := make([]GameCount, 0, 10)

	var gt string
	var count int
	for rows.Next() {
		err := rows.Scan(&gt, &count)
		if err != nil {
			log.Fatal("Error scanning rows for GetGameCounts()")
		}

		gc := GameCount{gt, count}
		gameCounts = append(gameCounts, gc)
	}

	return gameCounts
}

// GetSummaryStats gets high-level statistics for all-time or recently (with
// the "recent" flag). In particular it retrieves the number of active
// players and games played, but it also provides a breakdown of the top five
// most played game modes and their respective game counts. Game modes not in
// the top five are placed into an "other" category along with a count.
func GetSummaryStats(recent bool) SummaryStats {

	ds := SummaryStats{}

	// active players
	ds.Players = GetActivePlayers(recent)

	// note: this contains *all* game types, so we will have to check if it
	// needs to be pruned/condensed
	gcs := GetGameCounts(recent)

	games := 0
	otherGames := 0
	for i, gc := range gcs {
		games += gc.Games

		// condense other game types into an "other" category
		if i > 5 {
			otherGames += gc.Games
		}
	}
	ds.Games = games
	ds.OtherGames = otherGames

	if len(gcs) > 5 {
		ds.GameCounts = gcs[:5]
	} else {
		ds.GameCounts = gcs
	}

	return ds
}
