package riot_api

import (
	"fmt"
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	"strconv"
)

var client *golio.Client

var ApiKey string

// GetLolAPIClient returns a golio client
func GetLolAPIClient() *golio.Client {
	if client != nil {
		return client
	}
	client = golio.NewClient(ApiKey,
		golio.WithRegion(api.RegionEuropeWest))
	return client
}

// champCache is a cache for champion data, it is used to reduce the amount of requests to the riot api
var champCache = make(map[int]struct {
	Name    string
	Mastery *lol.ChampionMastery
})

// WipeCache wipes the champ mastery cache
func WipeCache() {
	champCache = make(map[int]struct {
		Name    string
		Mastery *lol.ChampionMastery
	})
}

// GetChampionMasteryById returns the mastery of a champion and there name
func GetChampionMasteryById(summonerName string, championId int) (*lol.ChampionMastery, string, error) {
	if champ, ok := champCache[championId]; ok {
		return champ.Mastery, champ.Name, nil
	}

	id, err := GetLolAPIClient().Riot.LoL.Summoner.GetByName(summonerName)
	if err != nil {
		fmt.Printf("Error getting summoner: %s\n", err)
		return nil, "", err
	}

	champion, err := GetLolAPIClient().DataDragon.GetChampionByID(strconv.Itoa(championId))
	if err != nil {
		fmt.Printf("Error getting champion: %s\n", err)
		return nil, "", err
	}

	fmt.Printf("Getting mastery for %s\n", champion.Name)
	mastery, err := GetLolAPIClient().Riot.LoL.ChampionMastery.Get(id.ID, champion.ID)
	if err != nil {
		fmt.Printf("Error getting mastery: %s\n", err)
		return nil, "", err
	}
	champCache[championId] = struct {
		Name    string
		Mastery *lol.ChampionMastery
	}{
		Name:    champion.Name,
		Mastery: mastery,
	}
	return mastery, champion.Name, err
}
