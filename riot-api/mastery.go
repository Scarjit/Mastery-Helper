package riot_api

import (
	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"log"
	"strconv"
)

var client *golio.Client

var ApiKey string

type Mastery struct {
	Name         string
	Level        int
	Points       int
	ChestGranted bool
	TokensEarned int
}

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
var champCache = make(map[int]Mastery)

// WipeCache wipes the champ mastery cache
func WipeCache() {
	champCache = make(map[int]Mastery)
}

// GetChampionMasteryById returns the mastery of a champion and there name
func GetChampionMasteryById(summonerName string, championId int) (*Mastery, error) {
	if champ, ok := champCache[championId]; ok {
		return &champ, nil
	}

	id, err := GetLolAPIClient().Riot.LoL.Summoner.GetByName(summonerName)
	if err != nil {
		log.Printf("Error getting summoner: (%s) %s\n", summonerName, err)
		return nil, err
	}

	champion, err := GetLolAPIClient().DataDragon.GetChampionByID(strconv.Itoa(championId))
	if err != nil {
		log.Printf("Error getting champion: (%d) %s\n", championId, err)
		return nil, err
	}

	mastery, err := GetLolAPIClient().Riot.LoL.ChampionMastery.Get(id.ID, champion.Key)
	m := Mastery{
		Name:         champion.Name,
		Level:        0,
		Points:       0,
		ChestGranted: false,
		TokensEarned: 0,
	}
	if err != nil {
		log.Printf("Error getting mastery for summoner %s and champion: %s: %s \n (Not played yet ?)", id.ID, champion.ID, err)
		champCache[championId] = m
		return &m, err
	}
	m = Mastery{
		Name:         champion.Name,
		Level:        mastery.ChampionLevel,
		Points:       mastery.ChampionPoints,
		ChestGranted: mastery.ChestGranted,
		TokensEarned: mastery.TokensEarned,
	}
	champCache[championId] = m
	return &m, err
}
