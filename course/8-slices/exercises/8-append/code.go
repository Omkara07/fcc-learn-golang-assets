package main

import "fmt"

type cost struct {
	day   int
	value float64
}

// simple solution
// func getCostsByDay(costs []cost) []float64 {
// 	ans := make([]float64, costs[len(costs)-1].day+1)
// 	for i := 0; i < len(costs); i++ {
// 		ans[costs[i].day] += costs[i].value
// 	}

// 	return ans
// }

// solution using append
func getCostsByDay(costs []cost) []float64 {
	// ans slice
	ans := []float64{}
	for i := 0; i < len(costs); i++ {
		// if the current day is >= the current length of the answer slice then expand the slice
		for costs[i].day >= len(ans) {
			ans = append(ans, 0.0)
		}
		// update the cost of the day
		ans[costs[i].day] += costs[i].value
	}

	return ans
}

// dont edit below this line

func test(costs []cost) {
	fmt.Printf("Creating daily buckets for %v costs...\n", len(costs))
	costsByDay := getCostsByDay(costs)
	fmt.Println("Costs by day:")
	for i := 0; i < len(costsByDay); i++ {
		fmt.Printf(" - Day %v: %.2f\n", i, costsByDay[i])
	}
	fmt.Println("===== END REPORT =====")
}

func main() {
	test([]cost{
		{0, 1.0},
		{1, 2.0},
		{1, 3.1},
		{2, 2.5},
		{3, 3.6},
		{3, 2.7},
		{4, 3.34},
	})
	test([]cost{
		{0, 1.0},
		{10, 2.0},
		{3, 3.1},
		{2, 2.5},
		{1, 3.6},
		{2, 2.7},
		{4, 56.34},
		{13, 2.34},
		{28, 1.34},
		{25, 2.34},
		{30, 4.34},
	})
}
