package controllers

import (
    "github.com/antzucaro/gostat/models"
    "github.com/antzucaro/gostat/templates"
    "html/template"
    "net/http"

    "fmt"
)

func Leaderboard(w http.ResponseWriter, r *http.Request) {
    type data struct {
        DailyStatLine template.HTML
        DuelRanks []models.PlayerRank
        CTFRanks []models.PlayerRank
        DMRanks []models.PlayerRank
    }
    var d data

    d.DuelRanks = models.GetTopNRanks("duel", 10)
    d.CTFRanks = models.GetTopNRanks("ctf", 10)
    d.DMRanks = models.GetTopNRanks("dm", 10)

    dsl := models.GetSummaryStats(true)

    dailyStatLine := fmt.Sprintf("%d active players and %d games (", 
        dsl.Players, dsl.Games)

    // common case is we have > 5 game modes to show
    if dsl.OtherGames > 0 {
        for _, gc := range dsl.GameCounts {
            dailyStatLine += fmt.Sprintf("%d %s, ", gc.Games, gc.GameType)
        }
        dailyStatLine += fmt.Sprintf("%d other", dsl.OtherGames)
    // less common is we append all that we have and don't show an "other"
    } else {
        // we have to construct the daily stat line from what we have thus far
        dailyStatLine += "TBD)"
    }
    dailyStatLine += fmt.Sprintf(") in the past 24 hours.")

    d.DailyStatLine = template.HTML(dailyStatLine)

    templates.Render("leaderboard", w, d)
}
