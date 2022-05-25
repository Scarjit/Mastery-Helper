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
	DisplayName string `json:"displayName"`
}
