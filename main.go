package main

import (
	"bytes"
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

func xkblayoutUpdater() chan string {
	c := make(chan string)

	update := func() {
		out, err := exec.Command("xkb-switch").Output()
		if err != nil {
			return
		}

		layout := string(out[:len(out)-1])
		c <- layout
	}

	go func() {
		update()
		for range newTickerTimer(time.Millisecond * 500) {
			update()
		}
	}()

	return c
}

func main() {
	updaters := [...]updater{
		bedTimeUpdater(),
		xkblayoutUpdater(),
		timeUpdater(),
	}

	cache := make(map[updater]string)

	var status bytes.Buffer

	for {
		select {
		case x := <-updaters[0]:
			cache[updaters[0]] = x
		case x := <-updaters[1]:
			cache[updaters[1]] = x
		case x := <-updaters[2]:
			cache[updaters[2]] = x
		}

		for _, c := range updaters {
			block := cache[c]
			if block == "" {
				continue
			}

			if status.Len() != 0 {
				status.WriteString(" | ")
			}

			status.WriteString(block)
		}

		exec.Command("xsetroot", "-name", status.String()).Run()

		status.Reset()
	}
}
