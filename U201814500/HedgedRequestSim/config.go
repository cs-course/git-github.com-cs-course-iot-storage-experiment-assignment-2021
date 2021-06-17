package main

const(
	lowDelay = 5 // 低延迟请求时间 ms
	highDelay = 1000 // 高延迟请求时间 ms
	highDelayRate = 0.2 // 请求遇到高延迟的概率
	maxTimeout = 10 // 从发送请求后允许的最大等待时间 ms
	candidates = 1 // 超过最大等待时间后下一次发出请求的数量
	testCount = 1000 // 测试次数
)


