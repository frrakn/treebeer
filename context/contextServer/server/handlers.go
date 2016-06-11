package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/frrakn/treebeer/context/db"
	"github.com/julienschmidt/httprouter"
)

func getPlayer(ps *players) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		paramID := p.ByName("id")

		id, err := strconv.Atoi(paramID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		player, ok := ps.get(db.PlayerID(id))
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		res, err := json.Marshal(player)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to render JSON for player with id: %s", paramID), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

func getTeam(ts *teams) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		paramID := p.ByName("id")

		id, err := strconv.Atoi(paramID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		team, ok := ts.get(db.TeamID(id))
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		res, err := json.Marshal(team)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to render JSON for team with id: %s", paramID), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

func getGame(gs *games) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		paramID := p.ByName("id")

		id, err := strconv.Atoi(paramID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		game, ok := gs.get(db.GameID(id))
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		res, err := json.Marshal(game)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to render JSON for game with id: %s", paramID), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}

func getStat(ss *stats) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		paramID := p.ByName("id")

		id, err := strconv.Atoi(paramID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		stat, ok := ss.get(db.StatID(id))
		if !ok {
			http.Error(w, fmt.Sprintf("Invalid id: %s", paramID), http.StatusBadRequest)
		}

		res, err := json.Marshal(stat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to render JSON for stat with id: %s", paramID), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}
