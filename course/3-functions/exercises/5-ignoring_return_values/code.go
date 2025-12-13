package main

import "fmt"

func main() {
	// ignoring an unused variable
	firstName, _ := getNames()
	fmt.Println("Welcome to Textio,", firstName)
}

// don't edit below this line

func getNames() (string, string) {
	return "John", "Doe"
}
