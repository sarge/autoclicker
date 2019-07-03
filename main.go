package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func main() {

	fmt.Println("Dads auto clicker")
	go listen()

	fmt.Println(" ` to start / stop")
	fmt.Println(" + to speed up the clicking")
	fmt.Println(" - to slow down the clicking")
	fmt.Println("ctrl + C to exit")
	select {}
}

func listen() {
	s := robotgo.Start()
	defer robotgo.End()

	var c chan int
	delay := 500

	stop := func() {
		close(c)
		c = nil
	}

	start := func(delay int) {
		c = clicker(delay) //ms
	}

	restart := func(delay int) {
		if c != nil {
			stop()
			start(delay)
		}
	}

	for ev := range s {
		// fmt.Println(ev)
		if ev.Kind == hook.KeyUp && ev.Rawcode == 192 {
			if c == nil {
				fmt.Println("started")
				start(delay)
			} else {
				fmt.Println("paused")
				stop()
			}
		}

		if ev.Kind == hook.KeyHold {
			switch ev.Rawcode {
			case 187: // speed up
				if delay > 100 {
					delay -= 100
					fmt.Printf("click delay %dms\n", delay)
					restart(delay)
				}
				break
			case 189: // slow down
				if delay < 5000 {
					delay += 100
					fmt.Printf("click delay %dms\n", delay)
					restart(delay)
				}
				break
			}
		}
	}
}

func clicker(delay int) chan int {
	ch := make(chan int)
	tick := time.Tick(time.Duration(delay) * time.Millisecond)

	go func() {
		for {
			select {
			case <-ch:
				return
			case <-tick:
				robotgo.Click()
			}
		}
	}()

	return ch
}
