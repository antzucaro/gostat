package models

import (
	"database/sql"
	"github.com/antzucaro/gostat/config"
	"github.com/antzucaro/gostat/qstr"
	"log"
	"strconv"
	"time"
)

// A PlayerTime represents the amount of time a given player has played during
// an unspecified time window.
type PlayerTime struct {
	N        int
	PlayerID int
	Nick     qstr.QStr
	PlayTime time.Duration
}

var playerTimeSQL = `select p.player_id, p.nick, 
extract(epoch from sum(pgs.alivetime)) playtime 
from players p join player_game_stats pgs 
on p.player_id = pgs.player_id 
and p.player_id > 2 
and pgs.create_dt > now() at time zone 'utc' - interval '` +
	config.Config.TopPlayersByTimeDays + ` days'
group by p.player_id, p.nick 
order by sum(pgs.alivetime) desc 
limit $1 
offset $2`

var playerTimeStmt *sql.Stmt

func GetTopPlayersByTime(limit int, offset int) []PlayerTime {
	rows, err := playerTimeStmt.Query(limit, offset)
	if err != nil {
		log.Fatal("Error processing GetTopPlayersByTime.")
	}

	playerTimes := make([]PlayerTime, 0, limit)

	n := 1
	var playerID int
	var nick string
	var playertime int

	for rows.Next() {
		rows.Scan(&playerID, &nick, &playertime)

		// convert the epoch from postgres into a duration-parseable string
		pts := strconv.Itoa(playertime) + "s"

		// now actually convert
		d, err := time.ParseDuration(pts)
		if err != nil {
			log.Fatal("Error converting the alivetime value in GetTopPlayersByTime.")
		}

		pt := PlayerTime{n, playerID, qstr.QStr(nick), d}

		playerTimes = append(playerTimes, pt)

		n += 1
	}

	return playerTimes
}
