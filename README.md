# Overstats
This simple CLI tool allows to see an Overwatch 2 player's most played heroes and their win rates in competitive mode.

## Usage
```bash
overstats (--heroesCount | -hc) <1-5> [--role | -r] <tank | damage | support> <battleTag>
```
### Flags
- `--role` or `-r`: The role to display. Must be one of `tank`, `damage` or `support`.
- _(Optional)_ `--heroesCount` or `-hc`: The number of heroes to display. Must be between 1 and 5.
### BattleTag
The BattleTag must be in the following format: `name#1234`. The name must be between 3 and 12 characters long and can only contain letters, numbers, underscores or hyphens. The number must be between 0 and 9999.

## Build
To build the project, you need to have [Go](https://golang.org/) installed. Then, simply run the following command:
```bash
go build
```
It will create an executable file named `overstats` in the current directory.

## TODO
- [X] Display the most played heroes and their win rates
- [ ] Display player information (number of games, rank and region)