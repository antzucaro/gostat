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
        TDMRanks []models.PlayerRank
        TopPlayersByTime []models.PlayerTime
        TopPlayersByScore []models.PlayerScore
        TopServersByGames []models.ServerGames
        RecentGames []models.RecentGame
    }

    // check the cache first
    cachedLeaderboard, ok := cache["Leaderboard"]
    if ok {
        println("Found in cache!")
        templates.Render("leaderboard", w, cachedLeaderboard)
        return
    }

    // otherwise construct it new
    var d data

    // ranks
    d.DuelRanks = models.GetTopNRanks("duel", 10)
    d.CTFRanks = models.GetTopNRanks("ctf", 10)
    d.DMRanks = models.GetTopNRanks("dm", 10)
    d.TDMRanks = models.GetTopNRanks("tdm", 10)

    // the overall stat line
    oss := models.GetSummaryStats(false)
    osl := makeStatLine("Tracking ", oss, " since October 2011.")
    d.OverallStatLine = osl

    // the daily stat line
    rss := models.GetSummaryStats(true)
    rsl := makeStatLine("", rss, " in the past 24 hours.")
    d.RecentStatLine = rsl

    // top players by playing time
    d.TopPlayersByTime = models.GetTopPlayersByTime(10, 0)

    // top players by accumulated score
    d.TopPlayersByScore = models.GetTopPlayersByScore(10, 0)

    // top servers by number of games
    d.TopServersByGames = models.GetTopServersByGames(10, 0)

    // recent games
    d.RecentGames = models.GetRecentGames(20, 0)

    // since we've done all this work, cache it!
    cache["Leaderboard"] = d

    templates.Render("leaderboard", w, d)
}

func makeStatLine(prefix string, stats models.SummaryStats, suffix string) template.HTML {
    // first check if there is anything to build
    if stats.Games == 0 {
        return ""
    }

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
