package main

import (
	"os/exec"
	"time"
)

type updater = chan string

func timeUpdater() chan string {
	const timeFormat = "Monday 02 15:04"

	c := make(chan string)

	update := func() { c <- time.Now().Format(timeFormat) }

	go func() {
		update()
		for range newTickerTimer(time.Minute) {
			update()
		}
	}()

	return c
}

func bedTimeUpdater() chan string {
	c := make(chan string)

	update := func() {
		now := time.Now()
		bedTime := time.Date(now.Year(), now.Month(), now.Day(), 22, 30, 0, 0, time.Local)

		if now.Before(bedTime) {
			c <- time.Until(bedTime).Truncate(time.Minute).String()
		} else {
			c <- "Bed Time!"
		}
	}

	go func() {
		update()
		for range newTickerTimer(time.Minute) {
			update()
		}
	}()

	return c
}

func main() {
	updaters := [...]updater{
		bedTimeUpdater(),
		timeUpdater(),
	}

	cache := make(map[updater]string)

	for {
		select {
		case x := <-updaters[0]:
			cache[updaters[0]] = x
		case x := <-updaters[1]:
			cache[updaters[1]] = x
		}

		status := ""
		for _, c := range updaters {
			block := cache[c]
			if block == "" {
				continue
			}

			if status != "" {
				status += " | "
			}

			status += block
		}

		exec.Command("xsetroot", "-name", status).Run()
	}
}
