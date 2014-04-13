package main

import (
    "github.com/antzucaro/gostat/controllers"
    "github.com/antzucaro/gostat/models"
    "github.com/antzucaro/gostat/templates"
    "github.com/go-martini/martini"
)

func main() {
  m := martini.Classic()

  // establish database connection, prepare queries
  models.Init("user=xonstat host=localhost dbname=xonstatdb sslmode=disable")

  // templates
  templates.Init()

  // routes
  m.Get("/", controllers.Leaderboard)

  m.Run()
}
