package main

import (
	"Mastery-Helper/helper"
	"Mastery-Helper/lcu"
	lol_champ_select "Mastery-Helper/lcu/lol-champ-select/v1"
	lol_login "Mastery-Helper/lcu/lol-summoner/v1/current-summoner"
	riot_api "Mastery-Helper/riot-api"
	"github.com/jedib0t/go-pretty/v6/table"
	"log"
	"os"
	"os/exec"
	"time"
)

var buildtime string

func main() {
	helper.InitLogging(buildtime)
	log.Printf("This is Mastery Helper v%s\n", buildtime)

	// Get the API key from the environment
	file, err := os.ReadFile("api_key.txt")
	if err != nil {
		log.Printf("Failed to read api_key.txt : %s\n", err)
		time.Sleep(time.Second * 5)
		return
	}
	riot_api.ApiKey = string(file)

	// Check if the API key is valid
	riotClient := riot_api.GetLolAPIClient()
	_, err = riotClient.Riot.LoL.Summoner.GetByName("Scarjit")
	if err != nil {
		log.Printf("Failed to retrieve scarjitSummoner, your API_KEY might be invalid : %s\n", err)
		return
	}

	for {
		riot_api.WipeCache()
		// Wait for league start
		lockfile := lcu.ReadLockFile()
		for lockfile == nil {
			time.Sleep(time.Second * 1)
			lockfile = lcu.ReadLockFile()
		}
		log.Printf("League started\n")

		// Get summoner
		summoner := lol_login.GetSummoner(lockfile.Password, lockfile.Port)
		for summoner == nil {
			time.Sleep(time.Second * 1)
			summoner = lol_login.GetSummoner(lockfile.Password, lockfile.Port)
		}
		log.Printf("Summoner logged in\n")

		// Wait for champ select
		session := lol_champ_select.GetSession(lockfile.Password, lockfile.Port)
		for session == nil {
			time.Sleep(time.Second * 1)
			log.Printf("Waiting for champ select\n")
			session = lol_champ_select.GetSession(lockfile.Password, lockfile.Port)
		}
		log.Printf("Champ select started\n")

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Champion", "Level", "Points", "Chest Granted", "Tokens Earned"})

		for session.Timer.AdjustedTimeLeftInPhase > 0 || session.Timer.IsInfinite {
			session = lol_champ_select.GetSession(lockfile.Password, lockfile.Port)
			if session == nil {
				log.Printf("Game was dodged !\n")
				break
			}

			var masteryInfo []*riot_api.Mastery
			if session.BenchEnabled {
				for _, championID := range session.BenchChampionIDS {
					championMastery, err := riot_api.GetChampionMasteryById(summoner.DisplayName, championID)
					if err != nil {
						log.Printf("Error while getting champion mastery : %s\n", err.Error())
						continue
					}
					masteryInfo = append(masteryInfo, championMastery)
				}
			}

			for _, team := range session.MyTeam {
				if team.ChampionID == 0 {
					continue
				}
				championMastery, err := riot_api.GetChampionMasteryById(summoner.DisplayName, int(team.ChampionID))
				if err != nil {
					log.Printf("Error getting mastery for %s\n", team.ChampionID)
					continue
				}
				masteryInfo = append(masteryInfo, championMastery)
			}

			t.ResetRows()
			for _, s := range masteryInfo {
				if s != nil {

					t.AppendRow(table.Row{
						s.Name,
						s.Level,
						s.Points,
						s.ChestGranted,
						s.TokensEarned,
					})
				}
			}

			ClearScreen()
			t.Render()

			time.Sleep(time.Second * 1)
		}
	}
}

func ClearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
