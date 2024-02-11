package poker_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/andremfp/poker-app"
)

// ScheduledAlert holds information about when an alert is scheduled.
type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

type SpyBlindAlerter struct {
	Alerts []ScheduledAlert
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, outputTo io.Writer) {
	s.Alerts = append(s.Alerts, ScheduledAlert{at, amount})
}

func TestGameStart(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		store := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		game := poker.NewTexasHoldem(store, blindAlerter)
		game.Start(5, io.Discard)

		tests := []ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		assertSchedulingTests(t, tests, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		store := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		game := poker.NewTexasHoldem(store, blindAlerter)
		game.Start(7, io.Discard)

		// requirement is:
		// the number of player determines the amount of time before the blind goes up
		// there is a base time of 5min
		// for every player, 1 minute is added
		// e.g. for 7 players, the blind goes up every 12min (5 + 7)
		tests := []ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		assertSchedulingTests(t, tests, blindAlerter)
	})
}

func TestGameFinish(t *testing.T) {
	t.Run("Andre wins", func(t *testing.T) {
		store := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(store, blindAlerter)
		game.Finish("Andre")

		assertPlayerWin(t, store, "Andre")
	})

	t.Run("Chris wins", func(t *testing.T) {
		store := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewTexasHoldem(store, blindAlerter)
		game.Finish("Chris")

		assertPlayerWin(t, store, "Chris")
	})
}

func assertSchedulingTests(t testing.TB, tests []ScheduledAlert, blindAlerter *SpyBlindAlerter) {
	for i, want := range tests {
		if len(blindAlerter.Alerts) <= i {
			t.Fatalf("alert %d was not scheduled, %v", i, blindAlerter.Alerts)
		}

		got := blindAlerter.Alerts[i]
		assertScheduledAt(t, got, want)
	}
}

func assertScheduledAt(t testing.TB, got, want ScheduledAlert) {
	t.Helper()
	if got.Amount != want.Amount {
		t.Errorf("got blind amount of %d, wanted %d", got.Amount, want.Amount)
	}
	if got.At != want.At {
		t.Errorf("got scheduled time of %v, wanted %v", got.At, want.At)
	}
}
