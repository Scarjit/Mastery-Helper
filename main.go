package main

import (
	"AramHelper/lcu"
	lol_champ_select "AramHelper/lcu/lol-champ-select/v1"
	lol_login "AramHelper/lcu/lol-summoner/v1/current-summoner"
	riot_api "AramHelper/riot-api"
	"fmt"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"os/exec"
	"time"
)

var buildtime string

func main() {
	fmt.Printf("This is Mastery Helper v%s\n", buildtime)

	// Get the API key from the environment
	file, err := os.ReadFile("api_key.txt")
	if err != nil {
		fmt.Printf("Failed to read api_key.txt : %s\n", err)
		time.Sleep(time.Second * 5)
		return
	}
	riot_api.ApiKey = string(file)

	// Check if the API key is valid
	riotClient := riot_api.GetLolAPIClient()
	_, err = riotClient.Riot.LoL.Summoner.GetByName("Scarjit")
	if err != nil {
		fmt.Printf("Failed to retrieve scarjitSummoner, your API_KEY might be invalid : %s\n", err)
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
		fmt.Printf("League started\n")

		// Get summoner
		summoner := lol_login.GetSummoner(lockfile.Password, lockfile.Port)
		for summoner == nil {
			time.Sleep(time.Second * 1)
			summoner = lol_login.GetSummoner(lockfile.Password, lockfile.Port)
		}
		fmt.Printf("Summoner logged in\n")

		// Wait for champ select
		session := lol_champ_select.GetSession(lockfile.Password, lockfile.Port)
		for session == nil {
			time.Sleep(time.Second * 1)
			session = lol_champ_select.GetSession(lockfile.Password, lockfile.Port)
		}
		fmt.Printf("Champ select started\n")

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Champion", "Level", "Points", "Chest Granted", "Tokens Earned"})

		for session.Timer.AdjustedTimeLeftInPhase > 0 {
			session = lol_champ_select.GetSession(lockfile.Password, lockfile.Port)
			if session == nil {
				fmt.Printf("Game was dodged !\n")
				break
			}

			var masteryInfo []struct {
				Name    string
				Mastery *lol.ChampionMastery
			}
			for _, championID := range session.BenchChampionIDS {
				championMastery, championName, err := riot_api.GetChampionMasteryById(summoner.DisplayName, championID)
				if err != nil {
					fmt.Printf("Error while getting champion mastery : %s\n", err.Error())
					continue
				}
				masteryInfo = append(masteryInfo, struct {
					Name    string
					Mastery *lol.ChampionMastery
				}{
					Name:    championName,
					Mastery: championMastery,
				})
			}

			for _, team := range session.MyTeam {
				if team.ChampionID == 0 {
					continue
				}
				championMastery, championName, err := riot_api.GetChampionMasteryById(summoner.DisplayName, int(team.ChampionID))
				if err != nil {
					fmt.Printf("Error getting mastery for %s\n", championName)
					continue
				}
				masteryInfo = append(masteryInfo, struct {
					Name    string
					Mastery *lol.ChampionMastery
				}{
					Name:    championName,
					Mastery: championMastery,
				})
			}

			t.ResetRows()
			for _, s := range masteryInfo {

				t.AppendRow(table.Row{
					s.Name,
					s.Mastery.ChampionLevel,
					s.Mastery.ChampionPoints,
					s.Mastery.ChestGranted,
					s.Mastery.TokensEarned,
				})
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
