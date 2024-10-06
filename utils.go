package main

import "time"

func newTickerTimer(updateInterval time.Duration) chan time.Time {
	c := make(chan time.Time)

	go func() {
		startTime := time.Now().Truncate(updateInterval).Add(updateInterval)
		start := time.NewTimer(startTime.Sub(time.Now()))

		timeUpdate := <-start.C
		c <- timeUpdate

		ticker := time.NewTicker(updateInterval)
		for timeUpdate = range ticker.C {
			c <- timeUpdate
		}
	}()

	return c
}
