package main

import (
	"fmt"
)

func printPrimes(max int) {
	for i := 2; i <= max; i++ {
		if i == 2 {
			fmt.Println(i)
			continue
		}
		if i%2 == 0 {
			continue
		}
		isPrime := true

		// for n := 3; n <= int(math.Sqrt(float64(i))); n++ {
		for n := 3; n*n <= i; n++ {

			if i%n == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			fmt.Println(i)
		}
	}
}

// don't edit below this line

func test(max int) {
	fmt.Printf("Primes up to %v:\n", max)
	printPrimes(max)
	fmt.Println("===============================================================")
}

func main() {
	test(10)
	test(20)
	test(30)
}
