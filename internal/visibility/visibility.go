package visibility

import (
	"log"
	"math"
	"sort"
	"time"
)

const epsilon = 1e-7 // Small value to avoid division by zero
const minObservableAltitude = 20.0
const maxObservableAltitude = 80.0

type VisibilityInfo struct {
	Object            AstroObject        `json:"object"`
	VisibilityWindows []VisibilityWindow `json:"visibilityWindows"`
}

type VisibilityWindow struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	StartAlt  float64   `json:"startAlt"`
	EndAlt    float64   `json:"endAlt"`
}

// TODO: validate that time is in UTC
func CalculateAltitudeVisibility(astroObjects *AstroObjectArray, configArray *ConfigArray, startTime, endTime time.Time, stepInMinutes time.Duration, printVisibleOnly bool) []VisibilityInfo {
	var allInfo []VisibilityInfo
	for _, astroObject := range astroObjects.Objects {
		allWindows := make([]VisibilityWindow, 0)
		for _, config := range configArray.Configs {
			visibilityWindows := make([]VisibilityWindow, 0)
			// Implement log output showing config info
			log.Printf("Config: %+v\n", config)
			if !ObjectNeverVisible(astroObject, &config) && ObjectEverInAzimuthWindow(astroObject, &config) {
				var lastVisibilityWindow *VisibilityWindow
				min, max := getTelescopeMinMaxAltitute(&config, config.DirectAzimuth)
				log.Printf("Telescope min altitude: %.2f°, max altitude: %.2f° at %f° azimuth\n", min, max, config.DirectAzimuth)
				log.Println("Start visibility calculation cycle")
				for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(stepInMinutes * time.Minute) {
					log.Println("Calculating visibility for time:", t.Format(time.RFC3339))
					alt, az := radecToAltAz(astroObject, &config.Position, t)
					// log.Printf("Altitude: %.2f°, Azimuth: %.2f°\n", alt, az)
					visible := isVisible(alt, az, &config)

					if lastVisibilityWindow != nil {
						lastVisibilityWindow.EndAlt = alt
						lastVisibilityWindow.EndTime = t
					}

					if visible {
						log.Printf("Object is visible at azimuth %.2f° and altitude %.2f°\n", az, alt)
						if lastVisibilityWindow == nil {
							lastVisibilityWindow = &VisibilityWindow{
								StartTime: t,
								StartAlt:  alt,
								EndAlt:    alt,
								EndTime:   t,
							}
						}
					} else if lastVisibilityWindow != nil {
						endVisibilityWindow(&lastVisibilityWindow, &visibilityWindows)
					}
				}
				if lastVisibilityWindow != nil {
					endVisibilityWindow(&lastVisibilityWindow, &visibilityWindows)
				}
			} else {
				log.Printf("Object %s is never visible from the given location and configuration.\n", astroObject.Name)
			}
			allWindows = append(allWindows, visibilityWindows...)
		}
		// Merge overlapping windows from all configs
		merged := mergeVisibilityWindows(allWindows)
		visibilityInfo := VisibilityInfo{
			Object:            astroObject,
			VisibilityWindows: merged,
		}
		allInfo = append(allInfo, visibilityInfo)
	}
	return allInfo
}

// mergeVisibilityWindows merges overlapping or adjacent visibility windows into a single list
func mergeVisibilityWindows(windows []VisibilityWindow) []VisibilityWindow {
	if len(windows) == 0 {
		return nil
	}
	// Sort by start time
	sort.Slice(windows, func(i, j int) bool {
		return windows[i].StartTime.Before(windows[j].StartTime)
	})
	merged := []VisibilityWindow{windows[0]}
	for i := 1; i < len(windows); i++ {
		last := &merged[len(merged)-1]
		curr := windows[i]
		if curr.StartTime.Before(last.EndTime) || curr.StartTime.Equal(last.EndTime) { // Overlapping or adjacent
			if curr.EndTime.After(last.EndTime) {
				last.EndTime = curr.EndTime
				last.EndAlt = curr.EndAlt
			}
		} else {
			merged = append(merged, curr)
		}
	}
	return merged
}

func ObjectNeverVisible(astroObject AstroObject, config *Config) bool {
	// Calculate the maximum altitude the object can reach
	declinationRad := Deg2rad(astroObject.Dec.toDegree())
	latitudeRad := Deg2rad(config.Position.Latitude)
	maxAltitude := Rad2deg(math.Asin(math.Sin(declinationRad)*math.Sin(latitudeRad) + math.Cos(declinationRad)*math.Cos(latitudeRad)))
	log.Printf("Maximum altitude of the object: %.2f°\n", maxAltitude)
	return maxAltitude < minObservableAltitude // Assuming 20 degrees is the minimum observable altitude
}

// Quick check if the object ever comes into the visible azimuth window
func ObjectEverInAzimuthWindow(astroObject AstroObject, config *Config) bool {
	// Quick check: does the object's rise/set azimuth range overlap with the visible window?
	if ObjectNeverVisible(astroObject, config) {
		return false
	}

	lat := Deg2rad(config.Position.Latitude)
	dec := Deg2rad(astroObject.Dec.toDegree())
	minAlt := Deg2rad(minObservableAltitude) // Minimum observable altitude

	// Calculate hour angle when object is at min observable altitude
	cosH := (math.Sin(minAlt) - math.Sin(lat)*math.Sin(dec)) / (math.Cos(lat) * math.Cos(dec))
	if cosH < -1 || cosH > 1 {
		// Object never reaches min observable altitude
		return false
	}
	// H := math.Acos(cosH) // Not used

	// Calculate azimuths at rise/set (when altitude = minAlt)
	// Azimuth at rise: A_rise = arccos((sin(dec) - sin(lat)*sin(minAlt)) / (cos(lat)*cos(minAlt)))
	// Azimuth at set: 360 - A_rise
	// sinA := math.Cos(dec)*math.Sin(H) / math.Cos(minAlt) // Not used
	cosA := (math.Sin(dec) - math.Sin(lat)*math.Sin(minAlt)) / (math.Cos(lat) * math.Cos(minAlt))
	// Clamp cosA to [-1, 1] to avoid NaN from floating-point errors
	clampedCosA := math.Max(-1, math.Min(1, cosA))
	azimuthRise := Rad2deg(math.Acos(clampedCosA))
	azimuthSet := 360.0 - azimuthRise

	// Normalize azimuths
	if azimuthRise > azimuthSet {
		azimuthRise, azimuthSet = azimuthSet, azimuthRise
	}

	left := config.LeftAzimuthLimit
	right := config.RightAzimuthLimit

	// Check for overlap between [A_rise, A_set] and [left, right]
	// Handle wrap-around
	visible := false
	if left < right {
		visible = !(azimuthSet < left || azimuthRise > right)
	} else {
		// Window wraps around 0
		visible = !(azimuthSet < left && azimuthRise > right)
	}
	return visible
}

func endVisibilityWindow(lastVisibilityWindow **VisibilityWindow, visibilityWindows *[]VisibilityWindow) {
	*visibilityWindows = append(*visibilityWindows, **lastVisibilityWindow)
	*lastVisibilityWindow = nil
}

func isVisible(objectAltitute float64, objectAzimuth float64, config *Config) bool {
	isAzimuthVisible := isClockwise(config.LeftAzimuthLimit, objectAzimuth) && isClockwise(objectAzimuth, config.RightAzimuthLimit)
	if !isAzimuthVisible {
		log.Printf("Not visible. Azimuth %.2f° is outside of limits [%.2f°, %.2f°]\n", objectAzimuth, config.LeftAzimuthLimit, config.RightAzimuthLimit)
		return false
	}
	alphaMin, alphaMax := getTelescopeMinMaxAltitute(config, objectAzimuth)
	log.Printf("Telescope min altitude: %.2f°, max altitude: %.2f° at %f° azimuth\n", alphaMin, alphaMax, objectAzimuth)
	if objectAltitute < minObservableAltitude || objectAltitute > maxObservableAltitude {
		log.Printf("Not visible. Altitude %.2f° is outside of limits [%.2f°, %.2f°]\n", objectAltitute, minObservableAltitude, maxObservableAltitude)
		return false
	}
	isAltitudeVisible := objectAltitute >= alphaMin && objectAltitute <= alphaMax
	if !isAltitudeVisible {
		log.Printf("Not visible. Altitude %.2f° is outside of limits [%.2f°, %.2f°]\n", objectAltitute, alphaMin, alphaMax)
		return false
	}
	return true
}

func isClockwise(leftAzimuth, rightAzimuth float64) bool {
	diff := math.Mod((rightAzimuth - leftAzimuth + 360), 360)
	log.Printf("Clockwise check: leftAzimuth=%.2f°, rightAzimuth=%.2f°, diff=%.2f°\n", leftAzimuth, rightAzimuth, diff)
	return diff > 0 && diff < 180
}

func getTelescopeMinMaxAltitute(config *Config, objectAzimuth float64) (float64, float64) {
	angleDiff := Deg2rad(math.Abs(objectAzimuth - config.DirectAzimuth))
	if angleDiff > 180 {
		angleDiff = 360 - angleDiff
	}
	alphaMin := altitudeAtAzimuthDiff(config.FenceHeight-config.TelescopeHeight, config.DistanceToFence, angleDiff)
	alphaMax := altitudeAtAzimuthDiff(config.WindowHeight+config.FenceHeight-config.TelescopeHeight, config.DistanceToFence, angleDiff)
	return Rad2deg(alphaMin), Rad2deg(alphaMax)
}

func altitudeAtAzimuthDiff(actualFenceHeight, distanceToFence, angleDiff float64) float64 {
	if actualFenceHeight <= 0 {
		actualFenceHeight = epsilon // Avoid division by zero
	}
	return math.Atan(actualFenceHeight * math.Cos(angleDiff) / distanceToFence)
}
