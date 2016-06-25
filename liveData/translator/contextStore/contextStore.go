package contextStore

type ContextStore struct {
	players *players
	teams   *teams
	games   *games
	stats   *stats
	poller  *RPCPoller
	ctxClient
}
