package schema

type RiotSeason struct {
	SeasonName  string
	SeasonSplit string
	ProTeams    []RiotTeam
	ProPlayers  []RiotPlayer
}

type RiotTeam struct {
	Id        int32
	RiotId    int32
	Name      string
	ShortName string
}

type RiotPlayer struct {
	Id        int32
	RiotId    int32
	Name      string
	ProTeamId int32
}
