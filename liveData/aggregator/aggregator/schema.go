package aggregator

import (
	pb "github.com/frrakn/treebeer/liveData/liveData"
)

type key struct {
	playerID int32
	teamID   int32
	gameID   int32
	statID   int32
}

type value struct {
	jsonValue string
	timestamp int32
}

func fromPB(proto *pb.Stat) (key, value) {
	return key{
			playerID: proto.Playerid,
			teamID:   proto.Teamid,
			gameID:   proto.Gameid,
			statID:   proto.Statid,
		},
		value{
			jsonValue: proto.Jsonvalue,
			timestamp: proto.Timestamp,
		}
}

func toPB(k key, v value) *pb.Stat {
	return &pb.Stat{
		Playerid:  k.playerID,
		Teamid:    k.teamID,
		Gameid:    k.gameID,
		Statid:    k.statID,
		Jsonvalue: v.jsonValue,
		Timestamp: v.timestamp,
	}
}
