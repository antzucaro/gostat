package main

import (
    "github.com/antzucaro/gostat/controllers"
    "github.com/antzucaro/gostat/models"
    "github.com/antzucaro/gostat/templates"
    "github.com/go-martini/martini"
    "net/http"
)

func main() {
  m := martini.Classic()

  // establish database connection, prepare queries
  models.Init("user=xonstat host=localhost dbname=xonstatdb sslmode=disable")

  // templates
  templates.Init()

  m.Get("/", func(w http.ResponseWriter, r *http.Request) {
      controllers.Leaderboard(w, r)
  })

  m.Run()
}
