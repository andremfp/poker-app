package poker_test

import (
	"os"
	"testing"

	"github.com/andremfp/poker-app"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("get league from reader", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Andre", "Wins": 20},
			{"Name": "Chris", "Wins": 10}]`)
		defer cleanDatabase()

		store, err := poker.NewFsPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()
		want := []poker.Player{
			{"Andre", 20},
			{"Chris", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get league sorted", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Andre", "Wins":10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFsPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()
		want := []poker.Player{
			{"Chris", 33},
			{"Andre", 10},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("get player score from reader", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Andre", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFsPlayerStore(database)
		assertNoError(t, err)

		got := store.GetPlayerScore("Andre")
		want := 10

		assertPlayerScore(t, got, want)
	})

	t.Run("record win for existing players", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Andre", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := poker.NewFsPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Andre")
		got := store.GetPlayerScore("Andre")
		want := 11

		assertPlayerScore(t, got, want)
	})

	t.Run("record win for new player", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, `[]`)
		defer cleanDatabase()

		store, err := poker.NewFsPlayerStore(database)
		assertNoError(t, err)

		store.RecordWin("Andre")
		got := store.GetPlayerScore("Andre")
		want := 1

		assertPlayerScore(t, got, want)
	})

	t.Run("works with empty file", func(t *testing.T) {

		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := poker.NewFsPlayerStore(database)
		assertNoError(t, err)
	})
}

// FsStore file just for testing with cleanup
func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
}

func assertPlayerScore(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got score of %d, wanted %d", got, want)
	}
}
