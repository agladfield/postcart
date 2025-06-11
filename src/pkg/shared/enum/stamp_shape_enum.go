package enum

type StampShapeEnum int8

const (
	StampShapeUnknown StampShapeEnum = iota
	StampShapeRect
	StampShapeRectClassic
	StampShapeCircle
	StampShapeCircleClassic
)

func (s StampShapeEnum) IsCircular() bool {
	switch s {
	case StampShapeCircle:
		fallthrough
	case StampShapeCircleClassic:
		return true
	default:
		return false
	}
}

// © Arthur Gladfield
