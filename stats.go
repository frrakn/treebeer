package treebeer

type StatType interface {
	id() int32
	aggregate(int32, int32) int32
}

type StatId struct {
	gameId   int32
	playerId int32
	statType StatType
}

//todo(frrakn): define various stat types offered by riot
