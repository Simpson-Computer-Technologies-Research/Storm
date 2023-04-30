package utils

import "math/rand"

// The RandomString() function is used to generate
// a random string with the provided length.
func RandomString(length int) string {
	// Define Variables
	var (
		// Characters to use in the random string
		chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// The result byte sloice
		res []byte = make([]byte, length)
	)
	// For each index in the result byte
	for i := range res {
		// Set it to a random character
		res[i] = chars[rand.Intn(len(chars))]
	}
	return string(res)
}
