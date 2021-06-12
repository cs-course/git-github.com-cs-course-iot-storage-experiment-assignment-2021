package main

import (
	"math/rand"
	"time"
)

func server(n int, c chan int) {
	var latency int

	for i := 0; i < n; i++{
		r := rand.Intn(100)
		if r < 100 * highDelayRate{
			latency = highDelay
		} else {
			latency = lowDelay
		}
		select{
		case <- time.After(time.Millisecond * time.Duration(latency)):
			c <- latency
		case <- time.After(time.Millisecond * maxTimeout):
			Chan := make(chan int, candidates)

			go server(candidates, Chan)
			select{
				case nextTime := <-Chan:
					c <- maxTimeout + nextTime
				case <-time.After(time.Millisecond * time.Duration(latency-maxTimeout)):
					c <- latency
			}
		}
	}
	close(c)
}


