package main

import (
	"fmt"
	"time"
)

// // 1st way
// func concurrrentFib(n int) {
// 	// create a channel
// 	fibCh := make(chan int)
// 	// recieve the numbers from the channel in a goroutine
// 	go func() {
// 		for num := range fibCh {
// 			fmt.Println(num)
// 		}
// 	}()
// 	// then call the fib func in the main thread
// 	fibonacci(n, fibCh)
// }

// 2nd way
func concurrrentFib(n int) {
	// create a channel
	fibCh := make(chan int)

	// call the fib func as a goroutine
	go fibonacci(n, fibCh)

	// then recieve the numbers in the main thread
	for num := range fibCh {
		fmt.Println(num)
	}
}

// TEST SUITE - Don't touch below this line

func test(n int) {
	fmt.Printf("Printing %v numbers...\n", n)
	concurrrentFib(n)
	fmt.Println("==============================")
}

func main() {
	test(10)
	test(5)
	test(20)
	test(13)
}

func fibonacci(n int, ch chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		ch <- x
		x, y = y, x+y
		time.Sleep(time.Millisecond * 10)
	}
	close(ch)
}
