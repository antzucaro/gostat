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
    GameCounts []GameCount
    OtherGames int
}

const dailyActivePlayersSQL = `SELECT count(distinct player_id) 
FROM player_game_stats 
WHERE player_id > 1
AND create_dt >= now() at time zone 'utc' - interval '20 days'`

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
WHERE create_dt >= now() at time zone 'utc' - interval '20 days'
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

// GetDailyStatLine gets high-level statistics from the past 24 hours. 
// In particular it retrieves the number of active players and games played,
// but it also provides a breakdown of the top five most played game modes
// and their respective game counts. Game modes not in the top five are 
// placed into an "other" category along with a count.
func GetDailyStatLine() DailyStat {
    ds := DailyStat{}

    // active players
    ds.Players = GetDailyActivePlayers()

    // note: this contains *all* game types, so we will have to check if it
    // needs to be pruned/condensed
    gcs := GetDailyGameCounts()

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
