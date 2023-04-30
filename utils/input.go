package utils

import "fmt"

// Function to get amount of followers to add
func InputInt(text string) int {
	var i int
	fmt.Print(text)
	fmt.Scanf("%d", &i)
	return i
}
