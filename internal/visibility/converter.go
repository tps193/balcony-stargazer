/*
Generated with AI
*/

package visibility

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/tps193/balcony-stargazer/internal/database"
)

// degrees to radians
func Deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// radians to degrees
func Rad2deg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// normalize angle to 0-360 deg
func normalize360(deg float64) float64 {
	for deg < 0 {
		deg += 360
	}
	for deg >= 360 {
		deg -= 360
	}
	return deg
}

// julian date from time.Time in UTC
func julianDate(t time.Time) float64 {
	unixTime := float64(t.UTC().Unix())
	jd := 2440587.5 + unixTime/86400.0
	return jd
}

// greenwich sidereal time at given julian date
func greenwichSiderealTime(jd float64) float64 {
	T := (jd - 2451545.0) / 36525.0
	GST := 280.46061837 + 360.98564736629*(jd-2451545.0) + 0.000387933*T*T - T*T*T/38710000.0
	return normalize360(GST)
}

// main conversion: RA, Dec, Lat, Lon, Time → Alt, Az
func radecToAltAz(astroObject AstroObject, position *Position, observationTime time.Time) (altDeg, azDeg float64) {
	utc := observationTime.UTC()
	jd := julianDate(utc)
	GST := greenwichSiderealTime(jd)

	// Local Sidereal Time in degrees
	LST_deg := normalize360(GST + position.Longitude)

	// Hour Angle in degrees: HA = LST - RA (all in degrees)
	HA_deg := normalize360(LST_deg - astroObject.Ra.toDegree())
	if HA_deg > 180 {
		HA_deg -= 360
	}

	// Convert to radians
	HA_rad := Deg2rad(HA_deg)
	dec_rad := Deg2rad(astroObject.Dec.toDegree())
	lat_rad := Deg2rad(position.Latitude)

	// Compute altitude
	sinAlt := math.Sin(dec_rad)*math.Sin(lat_rad) + math.Cos(dec_rad)*math.Cos(lat_rad)*math.Cos(HA_rad)
	alt_rad := math.Asin(sinAlt)

	// Compute azimuth
	cosAz := (math.Sin(dec_rad) - math.Sin(alt_rad)*math.Sin(lat_rad)) / (math.Cos(alt_rad) * math.Cos(lat_rad))
	// Clamp to [-1, 1] to avoid NaN due to floating point rounding
	if cosAz < -1 {
		cosAz = -1
	} else if cosAz > 1 {
		cosAz = 1
	}
	az_rad := math.Acos(cosAz)

	// Adjust azimuth based on HA
	if math.Sin(HA_rad) > 0 {
		az_rad = 2*math.Pi - az_rad
	}

	return Rad2deg(alt_rad), normalize360(Rad2deg(az_rad))
}

func ToAstroObjects(catalogRows []database.CatalogRow) (*AstroObjectArray, error) {
	astroObjects := &AstroObjectArray{}
	astroObjects.Objects = []AstroObject{}
	for _, obj := range catalogRows {
		name := obj.Name
		if obj.Commonnames != "" {
			name = fmt.Sprintf("%s (%s)", obj.Commonnames, obj.Name)
		}
		log.Println("Processing object:", name, "RA:", obj.RA, "Dec:", obj.Dec)
		ra, err := parseCatalogRA(obj.RA)
		if err != nil {
			fmt.Println("Error parsing RA:", err)
			return nil, err
		}
		dec, err := parseCatalogDec(obj.Dec)
		if err != nil {
			fmt.Println("Error parsing Dec:", err)
			return nil, err
		}
		astroObjects.Objects = append(astroObjects.Objects, AstroObject{
			Name: name,
			Ra:   ra,
			Dec:  dec,
		})
	}
	return astroObjects, nil
}

// parseCatalogRA parses a string like "10:08:28.10" into a RightAscension struct
func parseCatalogRA(s string) (RightAscension, error) {
	var ra RightAscension
	// convert floats to int where needed
	n, err := fmt.Sscanf(s, "%f:%f:%f", &ra.Hour, &ra.Min, &ra.Sec)
	if err != nil {
		return ra, fmt.Errorf("failed to parse RA: %w", err)
	}
	if n != 3 {
		return ra, fmt.Errorf("RA string should have format HH:MM:SS.ss, got: %s", s)
	}
	return ra, nil
}

// parseCatalogDec parses a string like "+12:18:23.0" or "-12:18:23.0" into a Declination struct
func parseCatalogDec(s string) (Declination, error) {
	var dec Declination
	sign := 1.0
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		if s[0] == '-' {
			sign = -1.0
		}
		s = s[1:]
	}
	var deg float64
	n, err := fmt.Sscanf(s, "%f:%f:%f", &deg, &dec.Min, &dec.Sec)
	if err != nil {
		return dec, fmt.Errorf("failed to parse Dec: %w", err)
	}
	if n != 3 {
		return dec, fmt.Errorf("Dec string should have format [+/-]DD:MM:SS.ss, got: %s", s)
	}
	dec.Degree = sign * float64(deg)
	return dec, nil
}
