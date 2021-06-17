package main

import (
	"sync"
)

const(
	lowLatency = 5 // 低延迟
	highLatency = 1000 // 高延迟
	highLatencyRate = 0.01 // 高延迟率
	queueCount = 10 // 任务队列个数
	concurrencyLimit = 10 // 每个任务队列可以并发执行的任务数量
	testCount = 1000 // 测试次数
	ifTiedRequest = false // 是否使用TiedRequest
)

var concurrentMap map[int]int
var lock sync.Mutex

func record(i interface{}, latency int){
	lock.Lock()
	defer lock.Unlock()

	concurrentMap[i.(int)] = latency
}

func checkTask(i interface{}) bool{
	lock.Lock()
	defer lock.Unlock()

	_, ok := concurrentMap[i.(int)]
	return ok
}