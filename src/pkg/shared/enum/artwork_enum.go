// Package enum handles enum cases used throughout the program
package enum

type ArtworkEnum int8

const (
	ArtworkUnknown ArtworkEnum = iota
	ArtworkAttachment
	ArtworkMountains
	ArtworkLakeside
	ArtworkIslands
	ArtworkCity
)

// © Arthur Gladfield
