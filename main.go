package main

import (
	"httpproject/service"
	"log"
	"net/http"
	"os"
	"time"

	klog "github.com/go-kit/kit/log"
	kitdatadog "github.com/go-kit/kit/metrics/dogstatsd"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

// func main() {
// 	server := &service.PlayerServer{
// 		Store: NewInMemoryPlayerStore(),
// 	}
// 	if err := http.ListenAndServe(":8080", server); err != nil {
// 		log.Fatalf("could not listen on port 8080 %v", err)
// 	}
// }

func main() {
	var svc service.PlayerStore
	svc = NewInMemoryPlayerStore()
	r := mux.NewRouter()
	logger := klog.NewJSONLogger(os.Stdout)
	dd := kitdatadog.New("playerscores", logger)

	requestCount := dd.NewCounter("request.count", 1)
	requestLatency := dd.NewTiming("request.latency", 1)
	countResult := dd.NewHistogram("result.count", 1)

	svc = &instrumentingMiddleware{requestCount, requestLatency, countResult, svc}
	getScoreEndpoint := makeGetScoreEndpoint(svc)
	{
		getScoreEndpoint = loggingMiddlware(logger)(getScoreEndpoint)
	}

	postWinEndpoint := makePostWinEndpoint(svc)
	{
		postWinEndpoint = loggingMiddlware(logger)(postWinEndpoint)
	}

	getScoreHandler := httptransport.NewServer(
		getScoreEndpoint,
		decodeGetScoreRequest,
		encodeGetScoreResponse,
		httptransport.ServerBefore(
			beforeIDExtractor,
			beforeMethodExtractor,
			beforePATHExtractor),
	)

	postWinHandler := httptransport.NewServer(
		postWinEndpoint,
		decodePostWinRequest,
		encodePostWinResponse,
		httptransport.ServerBefore(
			beforeIDExtractor,
			beforeMethodExtractor,
			beforePATHExtractor),
	)

	r.Methods("GET").PathPrefix("/players/").Handler(getScoreHandler)
	r.Methods("POST").PathPrefix("/players/").Handler(postWinHandler)

	// log.Fatal(http.ListenAndServe(":8080", nil))
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("running")
	log.Fatal(srv.ListenAndServe())
}

type InMemoryPlayerStore struct {
	scores map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.scores[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.scores[name]++
}
