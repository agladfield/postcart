// Package random contains randomization utlity methods
package random

import "math/rand"

// set seed

func FromSlice[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}
