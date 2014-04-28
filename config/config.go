package config

import (
    "encoding/json"
    "fmt"
    "os"
)

type config struct {
    // database connection string
    ConnStr string

	// how far back to retrieve the "summary stats" line at the top of the page
	SummaryStatsDays  string

	// how far back to retrieve the "top players by time" stats
	TopPlayersByTimeDays  string

	// how far back to retrieve the "top players by score" stats
	TopPlayersByScoreDays string

	// how far back to retrieve the "top servers by games" stats
	TopServersByGamesDays string

    // how far back to retrieve recent games
    RecentGamesDays string
}

var Path = "./config.json"
var Config = new(config)

func init() {
	// defaults
    Config.ConnStr = "user=xonstat host=localhost dbname=xonstatdb sslmode=disable"
    Config.SummaryStatsDays = "30"
    Config.TopPlayersByTimeDays = "30"
    Config.TopPlayersByScoreDays = "30"
    Config.TopServersByGamesDays = "30"
    Config.RecentGamesDays = "30"

    // set the config file path via environment variable
	if ecp := os.Getenv("GOSTAT_CONFIG"); ecp != "" {
		Path = ecp
	}

	file, err := os.Open(Path)
	if err != nil {
		if len(Path) > 1 {
			fmt.Printf("Error: could not read config file %s.\n", Path)
		}
		return
	}
	decoder := json.NewDecoder(file)

	// overwrite in-mem config with new values
	err = decoder.Decode(Config)
	if err != nil {
		fmt.Printf("Error decoding file %s\n%s\n", Path, err)
	}
}
