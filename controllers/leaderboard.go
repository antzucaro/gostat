package controllers

import (
    "fmt"
    "github.com/antzucaro/gostat/models"
    "github.com/antzucaro/gostat/templates"
    "html/template"
    "net/http"
)

func Leaderboard(w http.ResponseWriter, r *http.Request) {
    type data struct {
        OverallStatLine template.HTML
        RecentStatLine template.HTML
        DuelRanks []models.PlayerRank
        CTFRanks []models.PlayerRank
        DMRanks []models.PlayerRank
    }
    var d data

    d.DuelRanks = models.GetTopNRanks("duel", 10)
    d.CTFRanks = models.GetTopNRanks("ctf", 10)
    d.DMRanks = models.GetTopNRanks("dm", 10)

    oss := models.GetSummaryStats(false)
    osl := makeStatLine("Tracking ", oss, " since October 2011.")
    d.OverallStatLine = osl

    rss := models.GetSummaryStats(true)
    rsl := makeStatLine("", rss, " in the past 24 hours.")
    d.RecentStatLine = rsl

    templates.Render("leaderboard", w, d)
}

func makeStatLine(prefix string, stats models.SummaryStats, suffix string) template.HTML {
    line := fmt.Sprintf("%s%d players and %d games (", 
        prefix, stats.Players, stats.Games)

    // common case is we have > 5 game modes to show
    if stats.OtherGames > 0 {
        for _, gc := range stats.GameCounts {
            line += fmt.Sprintf("%d %s, ", gc.Games, gc.GameType)
        }
        line += fmt.Sprintf(" and %d other)", stats.OtherGames)
    // less common is we append all that we have and don't show an "other"
    } else {
        end := len(stats.GameCounts) - 1
        for i, gc := range stats.GameCounts {
            line += fmt.Sprintf("%d %s", gc.Games, gc.GameType)
            if i < end {
                line += ", "
            }
        }
        line += ")"
    }
    line += suffix

    return template.HTML(line)
}
