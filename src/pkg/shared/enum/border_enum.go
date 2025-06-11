package enum

type BorderEnum int8

const (
	BorderUnknown BorderEnum = iota
	BorderStandard
	BorderLines
	BorderCubes
	BorderStripes
	BorderPhoto
)

// © Arthur Gladfield
