package poker

import (
	"io"
	"time"
)

type TexasHoldem struct {
	blindAlerter BlindAlerter
	store        PlayerStore
}

func NewTexasHoldem(store PlayerStore, blindAlerter BlindAlerter) *TexasHoldem {
	return &TexasHoldem{
		blindAlerter: blindAlerter,
		store:        store,
	}
}

func (g *TexasHoldem) Start(numPlayers int, alertsDestination io.Writer) {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, blind := range blinds {
		g.blindAlerter.ScheduleAlertAt(blindTime, blind, alertsDestination)
		blindTime += time.Duration(5+numPlayers) * time.Minute
	}
}

func (g *TexasHoldem) Finish(winner string) {
	g.store.RecordWin(winner)
}
