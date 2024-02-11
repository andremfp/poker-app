package poker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andremfp/poker-app"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := poker.NewFsPlayerStore(database)
	assertNoError(t, err)

	server := mustMakePlayerServer(t, store, &SpyGame{})
	playerName := "Andre"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(playerName))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(playerName))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(playerName))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(playerName))
		assertResponseStatusCode(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetLeagueRequest())

		got := getLeagueFromResponse(t, response.Body)
		want := []poker.Player{
			{"Andre", 3},
		}
		assertResponseStatusCode(t, response.Code, http.StatusOK)
		assertLeague(t, got, want)
		assertContentType(t, response, poker.JsonContentType)
	})

}
