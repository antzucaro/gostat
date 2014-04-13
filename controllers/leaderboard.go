package controllers

import (
    "github.com/antzucaro/gostat/models"
    "github.com/antzucaro/gostat/templates"
    "net/http"
)

func Leaderboard(w http.ResponseWriter, r *http.Request) {
    type data struct {
        DuelRanks []models.PlayerRank
        CTFRanks []models.PlayerRank
        DMRanks []models.PlayerRank
    }
    var d data

    d.DuelRanks = models.GetTopNRanks("duel", 10)
    d.CTFRanks = models.GetTopNRanks("ctf", 10)
    d.DMRanks = models.GetTopNRanks("dm", 10)
    templates.Render("leaderboard", w, d)
}
