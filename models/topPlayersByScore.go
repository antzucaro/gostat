package models

import (
	"database/sql"
	"github.com/antzucaro/gostat/config"
	"github.com/antzucaro/gostat/qstr"
	"log"
)

// A PlayerScore represents the accumulated score of a player over
// an unspecified time window.
type PlayerScore struct {
	N        int
	PlayerID int
	Nick     qstr.QStr
	Score    int
}

var playerScoreSQL = `select player_id, nick, sum(score)
from player_game_stats
where create_dt > now() at time zone 'utc' - interval '` +
	config.Config.TopPlayersByScoreDays + ` days'
and player_id >= 2
and nick != 'Anonymous Player' 
group by player_id, nick
order by 3 desc
limit $1 
offset $2`

var playerScoreStmt *sql.Stmt

// GetTopPlayersByScore retrieves the top-scoring players for a duration
// of time (configurable via the "playScoreDays" constant). The data
// returned can be be constrained to a specific window within the result set
// via the "limit" and "offset" parameters.
func GetTopPlayersByScore(limit int, offset int) []PlayerScore {
	rows, err := playerScoreStmt.Query(limit, offset)
	if err != nil {
		log.Fatal("Error running query playerScoreStmt.")
	}

	playerScores := make([]PlayerScore, 0, limit)

	n := 1
	var playerID int
	var nick string
	var score int

	for rows.Next() {
		rows.Scan(&playerID, &nick, &score)

		ps := PlayerScore{n, playerID, qstr.QStr(nick), score}

		playerScores = append(playerScores, ps)

		n += 1
	}

	return playerScores
}
