package main

import (
	"fmt"

	rx "github.com/alecthomas/gorx"
)

func main() {
	observable := rx.FromStringArray([]string{"one", "two", "three", "four", "five", "two", "six"})
	fmt.Printf("All: %v\n", observable.ToArray())
	distinct := observable.Distinct()
	fmt.Printf("Distinct: %v\n", distinct.ToArray())
	fmt.Printf("Distinct: %v\n", distinct.ToArray())
	fmt.Printf("ElementAt(2): %v\n", observable.ElementAt(2).ToArray())
	fmt.Printf("Filter: %v\n", observable.Filter(func(s string) bool { return s != "three" }).ToArray())
	fmt.Printf("First: %v\n", observable.First().ToArray())
	fmt.Printf("Last: %v\n", observable.Last().ToArray())
	fmt.Printf("MapString: %v\n", rx.FromIntArray([]int{1, 2, 3, 4, 5}).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).ToArray())
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	fmt.Printf("Channel: %v\n", rx.FromIntChannel(ch).ToArray())
	fmt.Printf("IntArray -> Channel -> IntChannel -> Array: %v\n", rx.FromIntChannel(rx.FromIntArray([]int{1, 2, 3, 4}).ToChannel()).ToArray())
	fmt.Printf("SkipLast: %v\n", rx.FromIntArray([]int{1, 2, 3, 4, 5}).SkipLast(3).ToArray())
}
