package poker_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/andremfp/poker-app"
)

type SpyGame struct {
	StartCalled bool
	StartedWith int
	BlindAlert  []byte

	FinishedCalled bool
	FinishedWith   string
}

func (g *SpyGame) Start(numPlayers int, alertsDestination io.Writer) {
	g.StartCalled = true
	g.StartedWith = numPlayers
	alertsDestination.Write(g.BlindAlert)
}

func (g *SpyGame) Finish(winner string) {
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {

	t.Run("start game with 3 players and finish with 'Andre' as winner", func(t *testing.T) {
		game := &SpyGame{}
		stdout := &bytes.Buffer{}

		input := strings.NewReader("3\nAndre wins\n")

		cli := poker.NewCLI(input, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertStartCalledWith(t, game, 3)
		assertFinishCalledWith(t, game, "Andre")
	})

	t.Run("start game with 8 players and finish with 'Chris' as winner", func(t *testing.T) {
		game := &SpyGame{}
		stdout := &bytes.Buffer{}

		input := strings.NewReader("8\nChris wins\n")

		cli := poker.NewCLI(input, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
		assertStartCalledWith(t, game, 8)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("print error on bad number of players input", func(t *testing.T) {
		input := strings.NewReader("abc\n")
		stdout := &bytes.Buffer{}

		game := &SpyGame{}
		cli := poker.NewCLI(input, stdout, game)
		cli.PlayPoker()

		assertGameNotStarted(t, game)
		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.InvalidPlayerErrorPrompt)
	})

	t.Run("print error on bad winner input", func(t *testing.T) {
		game := &SpyGame{}
		stdout := &bytes.Buffer{}

		input := strings.NewReader("3\nNot a good input")

		cli := poker.NewCLI(input, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.InvalidWinnerErrorPrompt)
	})
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	got := stdout.String()
	want := strings.Join(messages, "")

	if got != want {
		t.Errorf("got prompt %q, want %q", got, want)
	}
}

func assertStartCalledWith(t testing.TB, game *SpyGame, want int) {
	t.Helper()

	// retry assertion
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartedWith == want
	})

	if !passed {
		t.Errorf("wanted Start called with %d, got %d", want, game.StartedWith)
	}
}

func assertFinishCalledWith(t testing.TB, game *SpyGame, want string) {
	t.Helper()

	// retry assertion
	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishedWith == want
	})

	if !passed {
		t.Errorf("wanted Finish called with %q, got %q", want, game.FinishedWith)
	}

}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func assertGameNotStarted(t testing.TB, game *SpyGame) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}
