package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) LocationExplore(url *string) (RespDeepLocations, error) {

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return RespDeepLocations{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RespDeepLocations{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RespDeepLocations{}, err
	}

	var exploreResp RespDeepLocations
	if err := json.Unmarshal(data, &exploreResp); err != nil {
		return RespDeepLocations{}, err
	}

	return exploreResp, nil
}
