package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "overstats",
		Usage: "Display Overwatch 2 useful stats of a player",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "heroesCount",
				Aliases: []string{"hc"},
				Value:   "3",
				Usage:   "Number of heroes to display",
			},
			&cli.StringFlag{
				Name:    "role",
				Aliases: []string{"r"},
				Usage:   "Role to display",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				return fmt.Errorf("Error: battleTag is required")
			}

			err := checkHeroesCount(ctx.Int("heroesCount"))
			if err != nil {
				return err
			}
			err = checkRole(ctx.String("role"))
			if err != nil {
				return err
			}

			battleTag := ctx.Args().First()
			data, err := fetchPlayerData(battleTag)

			if err != nil {
				log.Fatal(err)
			}

			if data["error"] != nil {
				return fmt.Errorf("Error: %s", data["error"])
			}

			if checkPrivateCareer(data) {
				return fmt.Errorf("Error: %s's career is private", battleTag)
			}

			role := ctx.String("role")

			username := fetchUsername(data)
			endorsement := fetchEndorsement(data)
			rankInfo := generateRankInfo(data, role)

			heroes := getMostPlayedHeroes(getCompetitiveStats(data), ctx.Int("heroesCount"), role)

			playerInfo := PlayerInfo{
				Name:        username,
				Endorsement: endorsement,
				Rank:        rankInfo,
				HeroInfo:    heroes,
			}

			displayPlayerInfo(playerInfo)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func checkHeroesCount(heroesCount int) error {
	if heroesCount < 1 || heroesCount > 5 {
		return fmt.Errorf("Error: heroesCount must be between 1 and 5")
	}

	return nil
}

func checkRole(role string) error {
	if role == "" {
		return fmt.Errorf("Error: role is required")
	}

	if role != "tank" && role != "damage" && role != "support" {
		return fmt.Errorf("Error: role must be tank, damage or support")
	}

	return nil
}

func displayPlayerInfo(playerInfo PlayerInfo) {
	fmt.Println("Player info:")
	fmt.Printf("- %s: %d endorsement level\n", playerInfo.Name, playerInfo.Endorsement)
	fmt.Printf("- %s: %s %d\n", playerInfo.Rank.Role, playerInfo.Rank.Division, playerInfo.Rank.Tier)
	fmt.Printf("- Season %d\n", playerInfo.Rank.Season)
	heroes := playerInfo.HeroInfo
	fmt.Println("Heroes:")
	for _, hero := range heroes {
		fmt.Printf("- %s: %d games played, %.2f%% win rate\n", hero.Name, hero.NumberOfGames, hero.WinPercentage)
	}
}
