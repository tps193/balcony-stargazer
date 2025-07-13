/*
Generated with AI
*/

package main

import (
	"math"
	"time"
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
	year, month, day := t.UTC().Date()
	hour := float64(t.UTC().Hour()) + float64(t.UTC().Minute())/60 + float64(t.UTC().Second())/3600

	if month <= 2 {
		year -= 1
		month += 12
	}
	A := year / 100
	B := 2 - A + A/4
	JD := float64(int(365.25*float64(year+4716))) + float64(int(30.6001*float64(month+1))) + float64(day) + float64(B) - 1524.5 + hour/24.0
	return JD
}

// greenwich sidereal time at given julian date
func greenwichSiderealTime(jd float64) float64 {
	T := (jd - 2451545.0) / 36525.0
	GST := 280.46061837 + 360.98564736629*(jd-2451545.0) + 0.000387933*T*T - T*T*T/38710000.0
	return normalize360(GST)
}

// main conversion: RA, Dec, Lat, Lon, Time → Alt, Az
func radecToAltAz(astroObject *AstroObject, position *Position, observationTime time.Time) (altDeg, azDeg float64) {
	jd := julianDate(observationTime)
	GST := greenwichSiderealTime(jd)
	LST := normalize360(GST+position.Longtitude) / 15.0 // in hours

	HA := normalize360((LST - astroObject.Ra.toDegree()) * 15) // hour angle in degrees
	if HA > 180 {
		HA -= 360
	}

	HA_rad := Deg2rad(HA)
	dec_rad := Deg2rad(astroObject.Dec.toDegree())
	lat_rad := Deg2rad(position.Latitude)

	sinAlt := math.Sin(dec_rad)*math.Sin(lat_rad) + math.Cos(dec_rad)*math.Cos(lat_rad)*math.Cos(HA_rad)
	alt := math.Asin(sinAlt)

	cosAz := (math.Sin(dec_rad) - math.Sin(alt)*math.Sin(lat_rad)) / (math.Cos(alt) * math.Cos(lat_rad))
	az := math.Acos(cosAz)

	// Adjust azimuth based on hour angle
	if math.Sin(HA_rad) > 0 {
		az = 2*math.Pi - az
	}

	return Rad2deg(alt), normalize360(Rad2deg(az))
}
