package database

type Filter struct {
	ObjectType        *string
	MinMagnitude      float64
	MinSizeArcMinutes float64
	MaxMagnitude      float64
	MaxSizeArcMinutes float64
}
