package current_summoner

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetSummoner returns summoner information of the currently logged in player
func GetSummoner(token string, port uint64) *Summoner {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://127.0.0.1:%d/lol-summoner/v1/current-summoner", port), nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("riot:%s", token)))))

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != 200 {
		return nil
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	session, err := UnmarshalSummoner(bytes)
	if err != nil {
		return nil
	}

	return &session
}

func UnmarshalSummoner(data []byte) (Summoner, error) {
	var r Summoner
	err := json.Unmarshal(data, &r)
	return r, err
}

type Summoner struct {
	AccountID                   int64        `json:"accountId"`
	DisplayName                 string       `json:"displayName"`
	InternalName                string       `json:"internalName"`
	NameChangeFlag              bool         `json:"nameChangeFlag"`
	PercentCompleteForNextLevel int64        `json:"percentCompleteForNextLevel"`
	Privacy                     string       `json:"privacy"`
	ProfileIconID               int64        `json:"profileIconId"`
	Puuid                       string       `json:"puuid"`
	RerollPoints                RerollPoints `json:"rerollPoints"`
	SummonerID                  int64        `json:"summonerId"`
	SummonerLevel               int64        `json:"summonerLevel"`
	Unnamed                     bool         `json:"unnamed"`
	XPSinceLastLevel            int64        `json:"xpSinceLastLevel"`
	XPUntilNextLevel            int64        `json:"xpUntilNextLevel"`
}

type RerollPoints struct {
	CurrentPoints    int64 `json:"currentPoints"`
	MaxRolls         int64 `json:"maxRolls"`
	NumberOfRolls    int64 `json:"numberOfRolls"`
	PointsCostToRoll int64 `json:"pointsCostToRoll"`
	PointsToReroll   int64 `json:"pointsToReroll"`
}
