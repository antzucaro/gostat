package models

import (
    "database/sql"
    "github.com/antzucaro/gostat/config"
    "github.com/antzucaro/gostat/qstr"
    "log"
    "time"
)

// A RecentGame represents a recently played game.
type RecentGame struct {
    GameID int
    GameType string
    ServerID int
    ServerName qstr.QStr
    MapID int
    MapName string
    GameDate time.Time
    WinnerID int
    WinnerNick qstr.QStr
}

var recentGamesSQL = `select g.game_id, g.game_type_cd, s.server_id, s.name, 
m.map_id, m.name, g.create_dt, pgs.player_id, pgs.nick
from games g join servers s on g.server_id = s.server_id
join maps m on g.map_id = m.map_id
join player_game_stats pgs on g.game_id = pgs.game_id
where pgs.scoreboardpos = 1
and g.create_dt > now() at time zone 'utc' - interval '` +
config.Config.RecentGamesDays + ` days'
order by g.game_id desc
limit $1
offset $2`

var recentGamesStmt *sql.Stmt

// GetRecentGames retrieves information for games played for a duration
// of time (configurable via the "recentGameDays" constant). The data
// returned can be be constrained to a specific window within the result set
// via the "limit" and "offset" parameters.
func GetRecentGames(limit int, offset int) []RecentGame {

    rows, err := recentGamesStmt.Query(limit, offset)
    if err != nil {
        log.Fatal("Error running query recentGamesStmt.")
    }

    recentGames := make([]RecentGame, 0, limit)

    var gameID int
    var gameType string
    var serverID int
    var serverName string
    var mapID int
    var mapName string
    var gameDate time.Time
    var winnerID int
    var winnerNick string

    for rows.Next() {
        rows.Scan(&gameID, &gameType, &serverID, &serverName, &mapID, &mapName, 
            &gameDate, &winnerID, &winnerNick)

        rg := RecentGame{gameID, gameType, serverID, qstr.QStr(serverName), 
            mapID, mapName, gameDate, winnerID, qstr.QStr(winnerNick)}

        recentGames = append(recentGames, rg)
    }

    return recentGames
}
