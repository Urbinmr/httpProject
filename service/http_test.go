package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}
func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestStoreWins(t *testing.T) {
	tests := []struct {
		name       string
		playerName string
		wantStatus int
		winCall    int
	}{
		{
			name:       "it returns accepted on POST",
			playerName: "Pepper",
			wantStatus: http.StatusAccepted,
			winCall:    1,
		},
		{
			name:       "it records wins when POST",
			playerName: "Pepper",
			wantStatus: http.StatusAccepted,
			winCall:    1,
		},
	}

	for _, test := range tests {
		store := StubPlayerStore{
			map[string]int{},
			[]string{},
		}
		server := &PlayerServer{&store}

		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodPost, "/players/"+test.playerName, nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)
			gotStatus := response.Code

			if gotStatus != test.wantStatus {
				t.Errorf("got status %d, want status %d", gotStatus, test.wantStatus)
			}
			if len(store.winCalls) != test.winCall {
				t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), test.winCall)
			}
			if store.winCalls[0] != test.playerName {
				t.Errorf("got name %s, want name %s", store.winCalls[0], test.playerName)
			}
		})
	}
}

func TestHTTPGet(t *testing.T) {
	tests := []struct {
		name       string
		playerName string
		want       string
		wantStatus int
	}{
		{
			name:       "returns Pepper's score",
			playerName: "Pepper",
			want:       "20",
			wantStatus: 200,
		},
		{
			name:       "returns Floyd's score",
			playerName: "Floyd",
			want:       "10",
			wantStatus: 200,
		},
		{
			name:       "returns a 404 on nonsense name",
			playerName: "Apple",
			want:       ErrPlayerNotFound.Error(),
			wantStatus: http.StatusNotFound,
		},
	}
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]string{},
	}
	server := &PlayerServer{&store}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/players/"+test.playerName, nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			gotBody := response.Body.String()
			gotStatus := response.Code

			if gotBody != test.want {
				t.Errorf("got '%s', want '%s'", gotBody, test.want)
			}
			if gotStatus != test.wantStatus {
				t.Errorf("got status %d, want status %d", gotStatus, test.wantStatus)
			}
		})
	}

}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name       string
		want       []Scores
		wantStatus int
	}{
		{
			name: "returns names and scores",
			want: []Scores{
				{
					name:  "Tim",
					score: 999,
				},
				{
					name:  "Matt",
					score: 4000,
				},
			},
			wantStatus: 200,
		},
	}
	store := StubPlayerStore{
		map[string]int{
			"Tim":  4000,
			"Matt": 999,
		},
		[]string{},
	}
	server := &PlayerServer{&store}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/scores/", nil)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			data, _ := ioutil.ReadAll(response.Body)
			var gotBody []Scores
			err := json.Unmarshal(data, &gotBody)
			if err != nil {
				panic(err)
			}
			gotStatus := response.Code

			if reflect.DeepEqual(gotBody, test.want) {
				t.Errorf("got '%v', want '%v'", gotBody, test.want)
			}
			if gotStatus != test.wantStatus {
				t.Errorf("got status %d, want status %d", gotStatus, test.wantStatus)
			}
		})
	}
}
