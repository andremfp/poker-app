package poker_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/andremfp/poker-app"
	"github.com/gorilla/websocket"
)

type StubPlayerStore struct {
	Scores   map[string]int
	WinCalls []string
	League   []poker.Player
}

func (s *StubPlayerStore) GetPlayerScore(playerName string) int {
	return s.Scores[playerName]
}

func (s *StubPlayerStore) RecordWin(playerName string) {
	s.WinCalls = append(s.WinCalls, playerName)
}

func (s *StubPlayerStore) GetLeague() poker.League {
	return s.League
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func newGetGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
}

func TestPlayerServer(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Andre": 20,
			"Chris": 40,
		},
		[]string{},
		nil,
	}
	server := mustMakePlayerServer(t, &store, &SpyGame{})

	t.Run("successful player get 1", func(t *testing.T) {
		request := newGetScoreRequest("Andre")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("successful player get 2", func(t *testing.T) {
		request := newGetScoreRequest("Chris")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "40")
	})

	t.Run("404 player not found", func(t *testing.T) {
		request := newGetScoreRequest("John")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		[]string{},
		nil,
	}
	server := mustMakePlayerServer(t, &store, &SpyGame{})

	t.Run("POST accepted", func(t *testing.T) {
		playerName := "Andre"
		request := newPostWinRequest(playerName)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusAccepted)
		assertPlayerWin(t, &store, playerName)
	})
}

func TestLeague(t *testing.T) {

	t.Run("/league returns 200", func(t *testing.T) {
		wantedLeague := []poker.Player{
			{"Andre", 32},
			{"Chris", 20},
			{"John", 13},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store, &SpyGame{})

		request := newGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, poker.JsonContentType)

	})
}

func TestGame(t *testing.T) {

	t.Run("get /game returns 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &StubPlayerStore{}, &SpyGame{})

		request := newGetGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("start a game with 3 players and declare Andre the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		winner := "Andre"

		game := &SpyGame{BlindAlert: []byte(wantedBlindAlert)}
		server := httptest.NewServer(mustMakePlayerServer(t, &StubPlayerStore{}, game))
		defer server.Close()

		// setup and test websocket connection
		// need a persistent connection to a server to test, hence the test server
		// calls /ws on the webserver to establish connection
		ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")
		defer ws.Close()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		assertStartCalledWith(t, game, 3)
		assertFinishCalledWith(t, game, winner)
		within(t, 10*time.Millisecond, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })

	})
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (got []poker.Player) {
	t.Helper()
	league, _ := poker.NewLeague(body)
	return league
}

func mustMakePlayerServer(t *testing.T, store poker.PlayerStore, game *SpyGame) *poker.PlayerServer {
	t.Helper()
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, wsURL string) *websocket.Conn {
	t.Helper()
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s, %v", wsURL, err)
	}
	return ws
}

func writeWSMessage(t *testing.T, wsConn *websocket.Conn, message string) {
	t.Helper()
	if err := wsConn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws, %v", err)
	}
}

func assertLeague(t testing.TB, got, want []poker.Player) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func assertResponseStatusCode(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got code %d, wanted %d", got, want)
	}
}

func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.WinCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin, want %d", len(store.WinCalls), 1)
	}

	if store.WinCalls[0] != winner {
		t.Errorf("did not get correct winner. got %q want %q", store.WinCalls[0], winner)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != "application/json" {
		t.Errorf("response does not have content-type header 'application/json', got %v", response.Result().Header)
	}
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()

	if string(msg) != want {
		t.Errorf("got %q, want %q", string(msg), want)
	}
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out waiting for message on ws")
	case <-done:
	}
}
