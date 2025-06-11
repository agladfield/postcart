// Package random contains randomization utlity methods
package random

import "math/rand"

// Create a local random source with the seed
var randSource *rand.Rand

func init() {
	SetSeed(rand.Int63())
}

func SetSeed(seed int64) {
	source := rand.NewSource(seed)
	randSource = rand.New(source)
}

func FromSlice[T any](slice []T) T {
	return slice[randSource.Intn(len(slice))]
}

// © Arthur Gladfield
