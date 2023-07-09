package main

type PlayerInfo struct {
	Name        string
	Endorsement int8
	Rank        RankInfo

	HeroInfo []HeroInfo
}

type RankInfo struct {
	Season   int
	Role     string
	Division string
	Tier     int
}

type HeroInfo struct {
	Name          string
	NumberOfGames int
	WinPercentage float64
}
