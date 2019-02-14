package main

import (
	"gopkg.in/logger.v1"
	"time"
)

func main() {
	for {
		log.Info("测试日志_1\n  测试日志_2\n  测试日志_3\n  测试日志_4\n  测试日志_5\n  测试日志_6\n  测试日志_7\n  测试日志_8")
		time.Sleep(time.Second * 1)
	}
}
