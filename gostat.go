package main

import (
    "github.com/antzucaro/gostat/db"
    "github.com/go-martini/martini"
    "html/template"
    "net/http"
)

func main() {
  m := martini.Classic()

  // establish database connection, prepare queries
  db.Init("user=xonstat host=localhost dbname=xonstatdb sslmode=disable")

  // templates
  main, err := template.ParseFiles("templates/base.html")
  if err != nil {
    panic(err)
  }

  m.Get("/", func(w http.ResponseWriter, r *http.Request) {
      type data struct {
          DuelRanks []db.PlayerRank
          CTFRanks []db.PlayerRank
          DMRanks []db.PlayerRank
      }
      var d data

      d.DuelRanks = db.GetTopNRanks("duel", 10)
      d.CTFRanks = db.GetTopNRanks("ctf", 10)
      d.DMRanks = db.GetTopNRanks("dm", 10)
      main.Execute(w, d)
  })

  m.Run()
}
