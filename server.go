package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/websocket"
)

const JsonContentType = "application/json"

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Player struct {
	Name string
	Wins int
}
type PlayerStore interface {
	GetPlayerScore(playerName string) int
	RecordWin(playerName string)
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
	game     Game
}

func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	// process and parse html template
	tmpl, err := template.ParseFiles("game.html")
	if err != nil {
		return nil, fmt.Errorf("problem loading template %s", err.Error())
	}

	p.game = game
	p.template = tmpl
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws", http.HandlerFunc(p.webSocketHandler))

	p.Handler = router

	return p, nil
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", JsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	playerName := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, r, playerName)
	case http.MethodGet:
		p.showScore(w, r, playerName)
	}
}

func (p *PlayerServer) gameHandler(w http.ResponseWriter, r *http.Request) {
	// write to w, meaning, display in the client (browser)
	p.template.Execute(w, nil)
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) webSocketHandler(w http.ResponseWriter, r *http.Request) {

	wsServer := NewPlayerServerWS(w, r)

	numberOfPlayersMsg := wsServer.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, wsServer)

	winner := wsServer.WaitForMsg()
	p.game.Finish(winner)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request, playerName string) {
	p.store.RecordWin(playerName)
	w.WriteHeader(http.StatusAccepted)

}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request, playerName string) {
	score := p.store.GetPlayerScore(playerName)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, score)
}
