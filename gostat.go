package main

import (
    _ "github.com/lib/pq"
    "database/sql"
    "fmt"
    "github.com/go-martini/martini"
)

// Player ranks
type PlayerRank struct {
    Id   int
    Nick string
    GameType string
    Elo  float64
    Rank int
}

const topNRanksSQL = "SELECT player_id, nick, game_type_cd, elo, rank FROM player_ranks WHERE game_type_cd = $1 AND rank <= $2"
var topNRanksStmt *sql.Stmt

func GetTopNRanks(db *sql.DB, gameType string, limit int) []PlayerRank {
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

        r := PlayerRank{id, nick, gameType, elo, rank}
        fmt.Println(r)
        fmt.Printf("%.0f\n", r.Elo)
        ranks = append(ranks, r)
    }
    err = rows.Err()
    if err != nil {
        panic(err)
    }

    return ranks
}

func main() {
  m := martini.Classic()

  // establish a database connection
  db, err := sql.Open("postgres", "user=xonstat host=localhost dbname=xonstatdb sslmode=disable")
  if err != nil {
    panic(err)
  }

  // prepare all of the queries
  topNRanksStmt, err = db.Prepare(topNRanksSQL)

  // all connections will have a db connection available
  m.Map(db)

  m.Get("/", func(db *sql.DB) string {
      GetTopNRanks(db, "duel", 10)
      GetTopNRanks(db, "ctf", 10)
      GetTopNRanks(db, "dm", 10)
      return "Hello, world!"
  })

  m.Run()
}
