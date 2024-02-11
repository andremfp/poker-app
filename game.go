package poker

import "io"

type Game interface {
	Start(numPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
