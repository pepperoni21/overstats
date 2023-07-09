package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func fetchPlayerData(battleTag string) (map[string]interface{}, error) {
	battleTag = formatBattleTag(battleTag)
	overFastAPIURL := fmt.Sprintf("https://overfast-api.tekrop.fr/players/%s", battleTag)

	response, err := http.Get(overFastAPIURL)

	if err != nil {
		return nil, err
	}

	return handleResponseAsMap(response)
}

func checkPrivateCareer(playerData map[string]interface{}) bool {
	summary := playerData["summary"].(map[string]interface{})
	privacy := summary["privacy"].(string)
	return privacy == "private"
}

func formatBattleTag(battleTag string) string {
	battleTag = strings.ReplaceAll(battleTag, "#", "-")
	battleTag = strings.TrimSpace(battleTag)
	return battleTag
}

func fetchHeroesByRole(role string) ([]string, error) {
	overFastAPIURL := fmt.Sprintf("https://overfast-api.tekrop.fr/heroes?role=%s", role)

	response, err := http.Get(overFastAPIURL)

	if err != nil {
		return nil, err
	}

	data, err := handleResponseAsList(response)

	if err != nil {
		return nil, err
	}

	heroes := make([]string, len(data))
	for i, hero := range data {
		hero := hero.(map[string]interface{})
		name := hero["key"].(string)
		heroes[i] = name
	}

	return heroes, nil
}

func handleResponseAsMap(response *http.Response) (map[string]interface{}, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func handleResponseAsList(response *http.Response) ([]interface{}, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
