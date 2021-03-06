package poller

import (
	"github.com/frrakn/treebeer/context/db"
	ctxPb "github.com/frrakn/treebeer/context/proto"
	"github.com/juju/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RPCPoller struct {
	*Poller
	client ctxPb.SeasonContextClient
}

func NewRPCPoller(address string) (*RPCPoller, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Trace(err)
	}

	p := &RPCPoller{
		client: ctxPb.NewSeasonContextClient(conn),
	}

	p.Poller = NewPoller(
		func() (*db.SeasonContext, error) {
			season := &db.SeasonContext{}
			fullCtx, err := p.client.GetContext(context.Background(), &ctxPb.Empty{})
			if err != nil {
				return nil, errors.Trace(err)
			}

			season.Players = make([]*db.Player, len(fullCtx.Players))
			for i, player := range fullCtx.Players {
				season.Players[i] = &db.Player{}
				err := season.Players[i].FromSavedPB(player)
				if err != nil {
					return nil, errors.Trace(err)
				}
			}

			season.Teams = make([]*db.Team, len(fullCtx.Teams))
			for i, team := range fullCtx.Teams {
				season.Teams[i] = &db.Team{}
				season.Teams[i].FromSavedPB(team)
			}

			season.Games = make([]*db.Game, len(fullCtx.Games))
			for i, game := range fullCtx.Games {
				season.Games[i] = &db.Game{}
				season.Games[i].FromSavedPB(game)
			}

			season.Stats = make([]*db.Stat, len(fullCtx.Stats))
			for i, stat := range fullCtx.Stats {
				season.Stats[i] = (&db.Stat{})
				season.Stats[i].FromSavedPB(stat)
			}

			return season, nil
		},
	)

	p.Poller.SetClose(func() {
		conn.Close()
	})

	return p, nil
}
