package main

import (
	"log"

	"golang.org/x/exp/slices"
)

func getCompetitiveStats(data map[string]interface{}) map[string]interface{} {
	statsSection := data["stats"].(map[string]interface{})
	pcSection := statsSection["pc"].(map[string]interface{})
	return pcSection["competitive"].(map[string]interface{})
}

func getMostPlayedHeroes(competitiveStats map[string]interface{}, heroesCount int, role string) []HeroInfo {
	mostPlayedHeroesMap := make(map[string]float64)

	heroesComparaison := competitiveStats["heroes_comparisons"].(map[string]interface{})
	timePlayed := heroesComparaison["time_played"].(map[string]interface{})["values"].([]interface{})

	heroes, err := fetchHeroesByRole(role)
	if err != nil {
		log.Fatal(err)
	}

	for _, hero := range timePlayed {
		hero := hero.(map[string]interface{})
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

func generateHeroInfo(heroName string, competitiveStats map[string]interface{}) HeroInfo {
	careerStats := competitiveStats["career_stats"].(map[string]interface{})
	heroCareerStats := careerStats[heroName].([]interface{})
	gameStats := heroCareerStats[3].(map[string]interface{})["stats"].([]interface{})

	gamesPlayed := gameStats[1].(map[string]interface{})["value"].(float64)
	gamesWon := gameStats[2].(map[string]interface{})["value"].(float64)

	winRate := gamesWon / gamesPlayed * 100

	return HeroInfo{
		Name:          heroName,
		NumberOfGames: int(gamesPlayed),
		WinPercentage: winRate,
	}
}
