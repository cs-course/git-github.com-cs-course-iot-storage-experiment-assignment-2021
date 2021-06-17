package main

import (
	"fmt"
	"math/rand"
	"time"
)

func task(val interface{}){
	r := rand.Intn(100)

	var latency int
	if r < 100 * highLatencyRate{
		latency = highLatency
	} else{
		latency = lowLatency
	}
	time.Sleep(time.Millisecond * time.Duration(latency))
	if !ifTiedRequest || !checkTask(val){
		record(val, latency)
	}
}

func main(){
	concurrentMap = make(map[int]int)

	queueGroup := make([]*Queue, queueCount)
	for i := 0; i < queueCount; i++{
		queueGroup[i] = NewQueue(task, concurrencyLimit)
	}

	for i := 0; i < testCount; i++ {
		if ifTiedRequest{
			for j := 0; j < queueCount; j++ {
				queueGroup[j].Push(i)
			}
		} else {
			index := rand.Intn(queueCount)
			queueGroup[index].Push(i)
		}
	}
	for i := 0; i < queueCount; i++ {
		queueGroup[i].Wait()
	}

	sum := 0
	for _, v := range concurrentMap{
		sum += v
	}
	fmt.Println("Avg Latency =", float64(sum) / testCount, "ms")
}
