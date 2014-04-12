package main

import (
    _ "github.com/lib/pq"
    "database/sql"
    "github.com/antzucaro/gostat/qstr"
    "github.com/go-martini/martini"
    "html/template"
    "net/http"
)

// Player ranks
type PlayerRank struct {
    Id   int
    Nick qstr.QStr
    GameType string
    Elo  float64
    Rank int
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

  // templates
  main, err := template.ParseFiles("templates/base.html")
  if err != nil {
    panic(err)
  }

  m.Get("/", func(w http.ResponseWriter, r *http.Request) {
      type data struct {
          DuelRanks []PlayerRank
          CTFRanks []PlayerRank
          DMRanks []PlayerRank
      }
      var d data

      d.DuelRanks = GetTopNRanks("duel", 10)
      d.CTFRanks = GetTopNRanks("ctf", 10)
      d.DMRanks = GetTopNRanks("dm", 10)
      main.Execute(w, d)
  })

  m.Run()
}
