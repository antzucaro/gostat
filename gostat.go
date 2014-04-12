package main

import (
    "github.com/antzucaro/gostat/models"
    "github.com/go-martini/martini"
    "html/template"
    "net/http"
)

func main() {
  m := martini.Classic()

  // establish database connection, prepare queries
  models.Init("user=xonstat host=localhost dbname=xonstatdb sslmode=disable")

  // templates
  main, err := template.ParseFiles("templates/base.html")
  if err != nil {
    panic(err)
  }

  m.Get("/", func(w http.ResponseWriter, r *http.Request) {
      type data struct {
          DuelRanks []models.PlayerRank
          CTFRanks []models.PlayerRank
          DMRanks []models.PlayerRank
      }
      var d data

      d.DuelRanks = models.GetTopNRanks("duel", 10)
      d.CTFRanks = models.GetTopNRanks("ctf", 10)
      d.DMRanks = models.GetTopNRanks("dm", 10)
      main.Execute(w, d)
  })

  m.Run()
}
