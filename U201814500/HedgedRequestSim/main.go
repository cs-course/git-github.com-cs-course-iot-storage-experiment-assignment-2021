package main

import "fmt"

func main() {
	c := make(chan int, testCount)
	go server(testCount, c)

	sum, cnt := 0, 0
	for elem := range c{
		sum += elem
		cnt ++
	}

	fmt.Println()
	fmt.Println("--- After", testCount, "tests:")
	fmt.Println("Avg Latency =", float64(sum)/ float64(cnt), "ms")
}
