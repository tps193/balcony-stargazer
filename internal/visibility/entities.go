package visibility

import "time"

const (
	VESPERA_HEIGHT = 18.00
)

type ConfigArray struct {
	Configs []Config `json:"configs"`
}

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

type AstroObjectArray struct {
	Objects []AstroObject `json:"objects"`
}

type ObjectType string

const (
	ObjectTypeStar             ObjectType = "*"
	ObjectTypeDoubleStar       ObjectType = "**"
	ObjectTypeStarAss          ObjectType = "*Ass"
	ObjectTypeClusterNebula    ObjectType = "Cl+N"
	ObjectTypeDuplicate        ObjectType = "Dup"
	ObjectTypeEmissionNebula   ObjectType = "EmN"
	ObjectTypeGalaxy           ObjectType = "G"
	ObjectTypeGalaxyCluster    ObjectType = "GCl"
	ObjectTypeGalaxyGroup      ObjectType = "GGroup"
	ObjectTypeGalaxyPair       ObjectType = "GPair"
	ObjectTypeGalaxyTriple     ObjectType = "GTrpl"
	ObjectTypeHIIRegion        ObjectType = "HII"
	ObjectTypeNebula           ObjectType = "Neb"
	ObjectTypeNonExistent      ObjectType = "NonEx"
	ObjectTypeNova             ObjectType = "Nova"
	ObjectTypeOpenCluster      ObjectType = "OCl"
	ObjectTypeOther            ObjectType = "Other"
	ObjectTypePlanetaryNebula  ObjectType = "PN"
	ObjectTypeReflectionNebula ObjectType = "RfN"
	ObjectTypeSupernovaRemnant ObjectType = "SNR"
)

type AstroObject struct {
	Name       string         `json:"name"`
	Ra         RightAscension `json:"ra"`
	Dec        Declination    `json:"dec"`
	ObjectType ObjectType     `json:"objectType"`
}

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RightAscension struct {
	Hour float64 `json:"hour"`
	Min  float64 `json:"min"`
	Sec  float64 `json:"sec"`
}

type Declination struct {
	Degree float64 `json:"degree"`
	Min    float64 `json:"min"`
	Sec    float64 `json:"sec"`
}

type TimeRange struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

type Filter struct {
	MinVisibilityDurationMinutes int `json:"minVisibilityDurationMinutes"`
	MinMagnitude                 int `json:"minMagnitude"`
	MinSizeArcMinutes            int `json:"minSizeArcMinutes"`
}

func (ra *RightAscension) toDegree() float64 {
	hours := float64(ra.Hour) + float64(ra.Min)/60.0 + float64(ra.Sec)/3600.0
	return hours * 15.0
}

func (dec *Declination) toDegree() float64 {
	return dec.Degree + dec.Min/60 + dec.Sec/3600
}
