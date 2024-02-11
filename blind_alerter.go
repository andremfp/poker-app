package poker

import (
	"fmt"
	"io"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, outputTo io.Writer)
}

type BlindAlerterFunc func(duration time.Duration, amount int, outputTo io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, outputTo io.Writer) {
	a(duration, amount, outputTo)
}

func Alerter(duration time.Duration, amount int, outputTo io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(outputTo, "Blind is now %d\n", amount)
	})
}
