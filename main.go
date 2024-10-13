package main

import "fmt"

func do(a int) {
	fmt.Println("Do")
	return
}

func first(a int) int {
	fmt.Println("First")
	return a + 1
}

func second() int {
	fmt.Println("Second")
	return 0
}

func domiddleware(a func(func())) {
	fmt.Println("Do Middleware")
	a(secondmiddleware)
}
func firstmiddleware(a func()) {
	fmt.Println("First Middleware")
	a()
}

func secondmiddleware() {
	fmt.Println("Second Middleware")
}

func main() {

	do(first(second()))
	domiddleware(firstmiddleware)
}
