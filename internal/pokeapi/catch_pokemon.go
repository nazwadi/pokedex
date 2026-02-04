package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) CatchPokemon(url *string) (Pokemon, error) {

	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var catchPokemon Pokemon
	if err := json.Unmarshal(data, &catchPokemon); err != nil {
		return Pokemon{}, err
	}

	return catchPokemon, nil
}
