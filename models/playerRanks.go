package models

import (
	"database/sql"
	"github.com/antzucaro/gostat/qstr"
)

// Player ranks
type PlayerRank struct {
	Id       int
	Nick     qstr.QStr
	GameType string
	Elo      float64
	Rank     int
}

const topNRanksSQL = "SELECT player_id, nick, game_type_cd, elo, rank FROM player_ranks WHERE game_type_cd = $1 AND rank <= $2"

var topNRanksStmt *sql.Stmt

func GetTopNRanks(gameType string, limit int) []PlayerRank {
	rows, err := topNRanksStmt.Query(gameType, limit)
	if err != nil {
		panic(err)
	}

	ranks := make([]PlayerRank, 0, 10)

	for rows.Next() {
		var id, rank int
		var nick, gameType string
		var elo float64

		err := rows.Scan(&id, &nick, &gameType, &elo, &rank)
		if err != nil {
			panic(err)
		}

		r := PlayerRank{id, qstr.QStr(nick), gameType, elo, rank}
		ranks = append(ranks, r)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return ranks
}
