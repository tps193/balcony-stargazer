package main

const (
	VESPERA_HEIGHT = 18.00
)

type Config struct {
	FenceHeight     float64
	WindowHeight    float64
	DistanceToFence float64
	TelescopeHeight float64
	DirectAzimuth   float64
	Position        Position
}

type AstroObject struct {
	Name string         `json:"name"`
	Ra   RightAscention `json:"ra"`
	Dec  Declanation    `json:"dec"`
}

type Position struct {
	Latitude   float64
	Longtitude float64
}

type RightAscention struct {
	Hour float64 `json:"hour"`
	Min  float64 `json:"min"`
	Sec  float64 `json:"sec"`
}

type Declanation struct {
	Degree float64 `json:"degree"`
	Min    float64 `json:"min"`
	Sec    float64 `json:"sec"`
}

func (ra RightAscention) toDegree() float64 {
	return ra.Hour + ra.Min/60 + ra.Sec/3600
}

func (dec Declanation) toDegree() float64 {
	return dec.Degree + dec.Min/60 + dec.Sec/3600
}
