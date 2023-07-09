package main

import (
	"log"

	"golang.org/x/exp/slices"
)

type JsonObject = map[string]interface{}

func getCompetitiveStats(playerData JsonObject, platform string) JsonObject {
	statsSection := playerData["stats"].(JsonObject)
	platformSection := statsSection[platform].(JsonObject)
	return platformSection["competitive"].(JsonObject)
}

func getMostPlayedHeroes(competitiveStats JsonObject, heroesCount int, role string) []HeroInfo {
	mostPlayedHeroesMap := make(map[string]float64)

	heroesComparaison := competitiveStats["heroes_comparisons"].(JsonObject)
	timePlayed := heroesComparaison["time_played"].(JsonObject)["values"].([]interface{})

	heroes, err := fetchHeroesByRole(role)
	if err != nil {
		log.Fatal(err)
	}

	for _, hero := range timePlayed {
		hero := hero.(JsonObject)
		heroName := hero["hero"].(string)

		if !slices.Contains(heroes, heroName) {
			continue
		}

		heroTimePlayed := hero["value"].(float64)

		mostPlayedHeroesMap[heroName] = heroTimePlayed
	}

	mostPlayedHeroesNames := make([]string, heroesCount)
	for i := 0; i < heroesCount; i++ {
		mostPlayedHero := getMostPlayedHero(mostPlayedHeroesMap)
		mostPlayedHeroesNames[i] = mostPlayedHero
		delete(mostPlayedHeroesMap, mostPlayedHero)
	}

	mostPlayedHeroes := make([]HeroInfo, heroesCount)
	for i, heroName := range mostPlayedHeroesNames {
		mostPlayedHeroes[i] = generateHeroInfo(heroName, competitiveStats)
	}

	return mostPlayedHeroes
}

func getMostPlayedHero(mostPlayedHeroesMap map[string]float64) string {
	currentHero := ""
	currentTimePlayed := 0.0

	for heroName, heroTimePlayed := range mostPlayedHeroesMap {
		if heroTimePlayed > currentTimePlayed {
			currentHero = heroName
			currentTimePlayed = heroTimePlayed
		}
	}

	return currentHero
}

func generateHeroInfo(heroName string, competitiveStats JsonObject) HeroInfo {
	careerStats := competitiveStats["career_stats"].(JsonObject)
	if careerStats[heroName] == nil {
		return HeroInfo{}
	}
	heroCareerStats := careerStats[heroName].([]interface{})
	gameStats := heroCareerStats[3].(JsonObject)["stats"].([]interface{})

	gamesPlayed := gameStats[1].(JsonObject)["value"].(float64)
	gamesWon := gameStats[2].(JsonObject)["value"].(float64)

	winRate := gamesWon / gamesPlayed * 100

	return HeroInfo{
		Name:          heroName,
		NumberOfGames: int(gamesPlayed),
		WinPercentage: winRate,
	}
}

func checkPrivateCareer(playerData JsonObject) bool {
	summary := getSummary(playerData)
	privacy := summary["privacy"].(string)
	return privacy == "private"
}

func getSummary(playerData JsonObject) JsonObject {
	return playerData["summary"].(JsonObject)
}

func generateRankInfo(playerData JsonObject, role string, platform string) RankInfo {
	summary := getSummary(playerData)
	competitive := summary["competitive"].(JsonObject)
	platformSection := competitive[platform].(JsonObject)
	if platformSection[role] == nil {
		return RankInfo{
			Role:     role,
			Season:   0,
			Tier:     0,
			Division: "",
		}
	}
	roleStats := platformSection[role].(JsonObject)

	season := platformSection["season"].(float64)
	division := roleStats["division"].(string)
	tier := roleStats["tier"].(float64)

	return RankInfo{
		Role:     role,
		Division: division,
		Tier:     int(tier),
		Season:   int(season),
	}
}

func fetchUsername(playerData JsonObject) string {
	summary := getSummary(playerData)
	username := summary["username"].(string)
	return username
}

func fetchEndorsement(playerData JsonObject) int8 {
	summary := getSummary(playerData)
	endorsement := summary["endorsement"].(JsonObject)
	level := endorsement["level"].(float64)
	return int8(level)
}
