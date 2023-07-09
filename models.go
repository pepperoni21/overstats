package main

type PlayerInfo struct {
	Name          string
	Region        string
	NumberOfGames int
	Rank          RankInfo

	HeroInfo []HeroInfo
}

type RankInfo struct {
	Role  string
	Rank  string
	Score int
}

type HeroInfo struct {
	Name          string
	NumberOfGames int
	WinPercentage float64
}
