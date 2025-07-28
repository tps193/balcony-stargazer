package main

import (
	"math"
	"time"
)

type VisibilityWindow struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	StartAlt  float64   `json:"startAlt"`
	EndAlt    float64   `json:"endAlt"`
}

func calculateAltitudeVisibility(astroObject *AstroObject, config *Config, startTime, endTime time.Time, stepInMinutes time.Duration, printVisibleOnly bool) []VisibilityWindow {
	visibilityWindows := make([]VisibilityWindow, 0)
	var lastVisibilityWindow *VisibilityWindow
	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(stepInMinutes * time.Minute) {

		alt, az := radecToAltAz(astroObject, &config.Position, t)
		visible := isVisible(alt, az, config)
		// t1, t2 := getTelescopeMinMaxAltitute(config, az)

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

		// if !printVisibleOnly || visible {
		// 	fmt.Println(t.Format("2006-01-02 15:04:05"))
		// 	// fmt.Printf("Altitude: %.2f°, Azimuth: %.2f°\n", alt, az)
		// 	fmt.Printf("Altitude: %.2f°\n", alt)
		// 	fmt.Println(visible)
		// 	fmt.Println(Rad2deg(t1), Rad2deg(t2))
		// }
	}
	if lastVisibilityWindow != nil {
		endVisibilityWindow(&lastVisibilityWindow, &visibilityWindows)
	}

	// fmt.Printf("Visibility of %s:\n", astroObject.Name)
	// for i, window := range visibilityWindows {
	// 	fmt.Printf("%d: %s\n", i, window.EndTime.Sub(window.StartTime))
	// 	fmt.Printf("\tStart: %s\n", window.StartTime)
	// 	fmt.Printf("\tEnd: %s\n", window.EndTime)
	// }
	return visibilityWindows
}

func endVisibilityWindow(lastVisibilityWindow **VisibilityWindow, visibilityWindows *[]VisibilityWindow) {
	*visibilityWindows = append(*visibilityWindows, **lastVisibilityWindow)
	*lastVisibilityWindow = nil
}

func isVisible(objectAltitute float64, objectAzimuth float64, config *Config) bool {
	alphaMin, alphaMax := getTelescopeMinMaxAltitute(config, objectAzimuth)
	return objectAltitute > alphaMin && objectAltitute < alphaMax
}

func getTelescopeMinMaxAltitute(config *Config, objectAzimuth float64) (float64, float64) {
	angleDiff := Deg2rad(math.Abs(objectAzimuth - config.DirectAzimuth))
	alphaMin := altitudeAtAzimuthDiff(config.FenceHeight-config.TelescopeHeight, config.DistanceToFence, angleDiff) - Deg2rad(3.0)
	alphaMax := altitudeAtAzimuthDiff(config.WindowHeight+config.FenceHeight-config.TelescopeHeight, config.DistanceToFence, angleDiff)
	return Rad2deg(alphaMin), Rad2deg(alphaMax)
}

func altitudeAtAzimuthDiff(actualFenceHeight, distanceToFence, angleDiff float64) float64 {
	return math.Atan(actualFenceHeight * math.Cos(angleDiff) / distanceToFence)
}
