package visibility

import (
	"log"
	"math"
	"time"
)

const epsilon = 1e-7 // Small value to avoid division by zero

type VisibilityWindow struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	StartAlt  float64   `json:"startAlt"`
	EndAlt    float64   `json:"endAlt"`
}

// TODO: validate that time is in UTC
func CalculateAltitudeVisibility(astroObject *AstroObject, config *Config, startTime, endTime time.Time, stepInMinutes time.Duration, printVisibleOnly bool) []VisibilityWindow {
	min, max := getTelescopeMinMaxAltitute(config, config.DirectAzimuth)
	log.Printf("Telescope min altitude: %.2f°, max altitude: %.2f° at %f° azimuth\n", min, max, config.DirectAzimuth)
	visibilityWindows := make([]VisibilityWindow, 0)
	var lastVisibilityWindow *VisibilityWindow
	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(stepInMinutes * time.Minute) {
		log.Println("Calculating visibility for time:", t.Format(time.RFC3339))
		alt, az := radecToAltAz(astroObject, &config.Position, t)
		log.Printf("Altitude: %.2f°, Azimuth: %.2f°\n", alt, az)
		visible := isVisible(alt, az, config)

		if lastVisibilityWindow != nil {
			lastVisibilityWindow.EndAlt = alt
			lastVisibilityWindow.EndTime = t
		}

		if visible {
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
	return visibilityWindows
}

func endVisibilityWindow(lastVisibilityWindow **VisibilityWindow, visibilityWindows *[]VisibilityWindow) {
	*visibilityWindows = append(*visibilityWindows, **lastVisibilityWindow)
	*lastVisibilityWindow = nil
}

func isVisible(objectAltitute float64, objectAzimuth float64, config *Config) bool {
	isAzimuthVisible := isClockwise(config.LeftAzimuthLimit, objectAzimuth) && isClockwise(objectAzimuth, config.RightAzimuthLimit)
	if !isAzimuthVisible {
		log.Printf("Object at azimuth %.2f° is not visible, outside of limits [%.2f°, %.2f°]\n", objectAzimuth, config.LeftAzimuthLimit, config.RightAzimuthLimit)
		return false
	}
	alphaMin, alphaMax := getTelescopeMinMaxAltitute(config, objectAzimuth)
	log.Printf("Telescope min altitude: %.2f°, max altitude: %.2f° at %f° azimuth\n", alphaMin, alphaMax, objectAzimuth)
	isAltitudeVisible := objectAltitute >= alphaMin && objectAltitute <= alphaMax
	if !isAltitudeVisible {
		log.Printf("Object at altitude %.2f° is not visible, outside of limits [%.2f°, %.2f°]\n", objectAltitute, alphaMin, alphaMax)
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
