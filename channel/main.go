package main

import "fmt"

var done = make(chan int, 2)

func count(nums []int, target int) {
	for i, num := range nums {
		for j, num1 := range nums {
			if num+num1 == target {
				done <- i
				done <- j
			}
		}
	}
}
func main() {
	nums := []int{2, 5, 3, 8, 0, 9, 23, 4}
	target := 12
	for i := 0; i <= 1; i++ {
		go count(nums, target)
		x := <-done
		y := <-done
		fmt.Println("x: ", x)
		fmt.Println("y: ", y)
	}
}

func Close(done chan bool) {
	close(done)
}
