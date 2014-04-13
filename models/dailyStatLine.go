package models

import (
    "database/sql"
    "log"
)

// The daily stat line is the summary line at the top of the page which shows
// the following info for the past 24 hours:
//
// - Number of active players
// - Number of games played
// - A breakdown of the number of games played by the top five game types  The
//   other game types played but not in the top five are grouped into "other".
type DailyStat struct {
    Players int
    Games int
    Type1 string
    Type1Games int
    Type2 string
    Type2Games int
    Type3 string
    Type3Games int
    Type4 string
    Type4Games int
    Type5 string
    Type5Games int
    OtherGames int
}

const dailyActivePlayersSQL = `SELECT count(distinct player_id) 
FROM player_game_stats 
WHERE player_id > 1
AND create_dt >= now() at time zone 'utc' - interval '1 day'`

var dailyActivePlayersStmt *sql.Stmt

// Retrieves the number of active players for the past 24 hours. Bots are 
// excluded, but // one anonymous player is included - we cannot count the 
// distinct number of those types of players (they are all player_id 2).
func GetDailyActivePlayers() int {
    row := dailyActivePlayersStmt.QueryRow()

    var dailyActivePlayers int
    err := row.Scan(&dailyActivePlayers)
    if err != nil {
        log.Fatal(err)
    }

    return dailyActivePlayers
}

// GameCount is a struct to hold the various game types and the
// number of times they have been played within a given interval
type GameCount struct {
    GameType string
    Games int
}

const dailyGamesSQL = `SELECT game_type_cd, count(*) 
FROM games 
WHERE create_dt >= now() at time zone 'utc' - interval '1 day'
GROUP BY game_type_cd
ORDER BY count(*) desc`

var dailyGamesStmt *sql.Stmt

// Retrieve the games played in the past 24 hours by game type.
// Returns an array of GameCount structs ordered from the most
// played game type to the least. No aggregations are performed.
func GetDailyGameCounts() []GameCount {
    rows, err := dailyGamesStmt.Query()
    if err != nil {
        log.Fatal(err)
    }

    gameCounts := make([]GameCount, 0, 10)

    var gt string
    var count int
    for rows.Next() {
        err := rows.Scan(&gt, &count)
        if err != nil {
            log.Fatal("Error scanning rows for GetDailyGameCounts()")
        }

        gc := GameCount{gt, count}
        gameCounts = append(gameCounts, gc)
    }

    return gameCounts
}
