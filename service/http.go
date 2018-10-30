package service

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrPlayerNotFound = errors.New("player not found")
)

type PlayerServer struct {
	Store PlayerStore
}

type PlayerStore interface {
	GetPlayerScore(string) int
	RecordWin(string)
	//postPlayerScore(string, int) error
}

type Scores struct {
	name  string
	score int
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		err := p.processWin(player)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusAccepted)
	case http.MethodGet:
		score, err := p.showScore(player)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		}
		if score != 0 {
			fmt.Fprintf(w, "%d", score)
			return
		}

	}
}

func (p *PlayerServer) processWin(player string) error {
	p.Store.RecordWin(player)
	return nil
}

func (p *PlayerServer) showScore(player string) (int, error) {
	score := p.Store.GetPlayerScore(player)
	if score == 0 {
		return score, ErrPlayerNotFound
	}
	return score, nil
}

func GetPlayerScore(player string) int {
	if player == "Pepper" {
		return 20
	}

	if player == "Floyd" {
		return 10
	}

	return 0
}

// func PostPlayerScore(name string, score int, ps PlayerStore) error {
// 	return ps.postPlayerScore(name, score)
// }
