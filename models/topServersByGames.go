package models

import (
    "database/sql"
    "github.com/antzucaro/gostat/qstr"
    "log"
)

// A ServerTime represents the number of games hosted by a server during
// an unspecified time window.
type ServerGames struct {
    ServerID int
    Name qstr.QStr
    Games int
}

const serverGamesDays = "7"

const serverGamesSQL = `select s.server_id, s.name, count(*) games
from servers s join games g on s.server_id = g.server_id
where g.create_dt > now() at time zone 'utc' - interval '` + 
serverGamesDays + ` days'
group by s.server_id, s.name
order by 3 desc
limit $1 
offset $2`

var serverGamesStmt *sql.Stmt

// GetTopServersByGames retrieves the servers playing the most games within an
// unspecified time period (as controlled by the "serverGamesDays" constant) in
// descending order. The limit and offset parameters can be used to move a
// limited-size window around within the result set.
func GetTopServersByGames(limit int, offset int) []ServerGames {
    rows, err := serverGamesStmt.Query(limit, offset)
    if err != nil {
        log.Fatal("Error running query serverGamesStmt.")
    }

    serverGames := make([]ServerGames, 0, limit)

    var serverID int
    var name string
    var games int

    for rows.Next() {
        rows.Scan(&serverID, &name, &games)

        sg := ServerGames{serverID, qstr.QStr(name), games}

        serverGames = append(serverGames, sg)
    }

    return serverGames
}
