package visibility

const (
	VESPERA_HEIGHT = 18.00
)

type Config struct {
	FenceHeight       float64  `json:"fenceHeight"`
	WindowHeight      float64  `json:"windowHeight"`
	DistanceToFence   float64  `json:"distanceToFence"`
	TelescopeHeight   float64  `json:"telescopeHeight"`
	DirectAzimuth     float64  `json:"directAzimuth"`
	Position          Position `json:"position"`
	LeftAzimuthLimit  float64  `json:"leftAzimuthLimit"`
	RightAzimuthLimit float64  `json:"rightAzimuthLimit"`
}

type AstroObject struct {
	Name string         `json:"name"`
	Ra   RightAscention `json:"ra"`
	Dec  Declanation    `json:"dec"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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

func (ra *RightAscention) toDegree() float64 {
	hours := float64(ra.Hour) + float64(ra.Min)/60.0 + float64(ra.Sec)/3600.0
	return hours * 15.0
}

func (dec *Declanation) toDegree() float64 {
	return dec.Degree + dec.Min/60 + dec.Sec/3600
}
