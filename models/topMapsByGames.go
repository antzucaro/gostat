package models

import (
    "database/sql"
    "github.com/antzucaro/gostat/config"
    "log"
)

// MapGames represents the number of games played on a map during
// an unspecified time window.
type MapGames struct {
    N int
    MapID int
    Name string
    Games int
}

var mapGamesSQL = `select m.map_id, m.name, count(*) games
from maps m join games g on m.map_id = g.server_id
where g.create_dt > now() at time zone 'utc' - interval '` + 
config.Config.TopMapsByGamesDays + ` days'
group by m.map_id, m.name
order by 3 desc
limit $1 
offset $2`

var mapGamesStmt *sql.Stmt

// GetTopMapsByGames retrieves the maps played the most within an
// unspecified time period (as controlled by the "mapGamesDays" constant) in
// descending order. The limit and offset parameters can be used to move a
// limited-size window around within the result set.
func GetTopMapsByGames(limit int, offset int) []MapGames {
    rows, err := mapGamesStmt.Query(limit, offset)
    if err != nil {
        log.Fatal("Error running query mapGamesStmt.")
    }

    mapGames := make([]MapGames, 0, limit)

    n := 1
    var mapID int
    var name string
    var games int

    for rows.Next() {
        rows.Scan(&mapID, &name, &games)

        mg := MapGames{n, mapID, name, games}

        mapGames = append(mapGames, mg)

        n += 1
    }

    return mapGames
}
