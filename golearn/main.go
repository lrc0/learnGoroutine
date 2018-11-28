package main

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	wg sync.WaitGroup
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"10.129.xx.12:9092", "10.129.xx.13:9092", "10.129.xx.14:9092"}, nil)
	if err != nil {
		panic(err)
	}

	partitionList, err := consumer.Partitions("testGo")

	if err != nil {
		panic(err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("topic", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}

		defer pc.AsyncClose()

		wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
		wg.Wait()
		consumer.Close()
	}
}

// package main

// import (
// 	"fmt"
// 	"time"
// )

// func main() {
// 	begin := "2018-07-16 00:00:00"
// 	end := "2018-07-16 08:00:00"
// 	// dt := time.Now().Format("2006-01-02 15:04:05")
// 	// hour, _ := time.ParseDuration("1h")
// 	// df := time.Now().Add(-hour).Format("2006-01-02 15:04:05")

// 	// fmt.Println("from: ", df)
// 	// fmt.Println("to: ", dt)
// 	slice := dealTime(begin, end)
// 	fmt.Println("slice: ", slice)
// }

// func dealTime(begin, end string) []string {
// 	timeLayout := "2006-01-02 15:04:05"
// 	timeFormat := "2006.01.02"
// 	loc, _ := time.LoadLocation("Local")

// 	beginTime, _ := time.ParseInLocation(timeLayout, begin, loc)
// 	endTime, _ := time.ParseInLocation(timeLayout, end, loc)
// 	nowTime := time.Now().Format(timeLayout)
// 	now, _ := time.ParseInLocation(timeLayout, nowTime, loc)
// 	if now.Before(endTime) {
// 		fmt.Println("=====================> now: ", now.String())
// 		fmt.Println("=====================> endtime: ", endTime.String())
// 		endTime = now
// 		fmt.Println("endtime: ", endTime)
// 	}

// 	beginTime = beginTime.Add(-8 * time.Hour)
// 	endTime = endTime.Add(-8 * time.Hour)

// 	begins := beginTime.Format(timeFormat)
// 	ends := endTime.Format(timeFormat)

// 	b, _ := time.Parse(timeFormat, begins)
// 	e, _ := time.Parse(timeFormat, ends)
// 	x := e.Sub(b).Hours() / 24

// 	var slice []string
// 	slice = append(slice, "logstash-"+begins)

// 	for i := 0; i < int(x); i++ {
// 		day, _ := time.ParseDuration("24h")
// 		b = b.Add(day)
// 		bs := b.Format(timeFormat)
// 		slice = append(slice, "logstash-"+bs)
// 	}
// 	return slice
// }
