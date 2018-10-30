package main

import (
	"context"
	"fmt"
	"httpproject/service"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func makeGetScoreEndpoint(svc service.PlayerStore) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getScoreRequest)
		v := svc.GetPlayerScore(req.player)
		if v == 0 {
			return getScoreResponse{v, service.ErrPlayerNotFound, http.StatusNotFound}, nil
		}
		return getScoreResponse{v, nil, http.StatusOK}, nil
	}
}

func makePostWinEndpoint(svc service.PlayerStore) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(postWinRequest)
		svc.RecordWin(req.player)
		return postWinResponse{http.StatusAccepted}, nil
	}
}

func decodeGetScoreRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getScoreRequest
	player := r.URL.Path[len("/players/"):]
	request.player = player
	return request, nil
}

func decodePostWinRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request postWinRequest
	player := r.URL.Path[len("/players/"):]
	request.player = player
	return request, nil
}

func encodePostWinResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// resp := response.(postWinResponse)
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func encodeGetScoreResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(getScoreResponse)
	if resp.err != nil {
		return resp.err
	}
	fmt.Fprintf(w, "%d", resp.score)
	return nil
}

type getScoreRequest struct {
	player string
}

type getScoreResponse struct {
	score int
	err   error
	code  int
}

func (r getScoreResponse) String() string {
	return fmt.Sprintf("Score: %d Status Code: %d", r.score, r.code)
}

type postWinRequest struct {
	player string
}

type postWinResponse struct {
	code int
}

func (r postWinResponse) String() string {
	return fmt.Sprintf("Status Code: %d", r.code)
}
