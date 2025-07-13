package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	config := Config{
		FenceHeight:     43.25,
		WindowHeight:    62.00,
		DistanceToFence: 35,
		TelescopeHeight: VESPERA_HEIGHT,
		DirectAzimuth:   80.00,
		Position:        Position{37.38, -121.89},
	}

	if len(os.Args) < 2 {
		fmt.Println("No object specified")
		return
	}

	jsonStr := os.Args[1]
	var object AstroObject
	err := json.Unmarshal([]byte(jsonStr), &object)
	if err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}

	// object := AstroObject{
	// 	Name: "NGC 6992",
	// 	Ra:   RightAscention{20, 58, 18},
	// 	Dec:  Declanation{31, 43, 0},
	// }

	startTime := time.Date(time.Now().Year(), 6, 30, 22, 30, 0, 0, time.Now().Location())
	endTime := time.Date(2025, startTime.Month(), startTime.Day()+1, 1, 30, 0, 0, time.Now().Location())

	fmt.Printf("Observed time from %s to %s\n", startTime, endTime)
	fmt.Println(config)
	fmt.Println(object)

	calculateAltitudeVisibility(&object, &config, startTime, endTime, 5, true)
}

type VisibilityWindow struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	StartAlt  float64   `json:"startAlt"`
	EndAlt    float64   `json:"endAlt"`
}

func calculateAltitudeVisibility(astroObject *AstroObject, config *Config, startTime, endTime time.Time, stepInMinutes time.Duration, printVisibleOnly bool) {
	visibilityWindows := make([]VisibilityWindow, 0)
	var lastVisibilityWindow *VisibilityWindow
	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(stepInMinutes * time.Minute) {

		alt, az := radecToAltAz(astroObject, &config.Position, t)
		visible := isVisible(alt, az, config)
		// t1, t2 := getTelescopeMinMaxAltitute(config, az)

		if visible {
			if lastVisibilityWindow == nil {
				lastVisibilityWindow = &VisibilityWindow{
					StartTime: t,
					StartAlt:  alt,
				}
			}
		} else if lastVisibilityWindow != nil {
			endVisibilityWindow(&lastVisibilityWindow, &visibilityWindows, t)
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
		endVisibilityWindow(&lastVisibilityWindow, &visibilityWindows, endTime)
	}

	var result Result
	result = NewJsonOutput() //NewSimpleOutputResult()
	fmt.Println(result.Get(astroObject, &visibilityWindows))

	// fmt.Printf("Visibility of %s:\n", astroObject.Name)
	// for i, window := range visibilityWindows {
	// 	fmt.Printf("%d: %s\n", i, window.EndTime.Sub(window.StartTime))
	// 	fmt.Printf("\tStart: %s\n", window.StartTime)
	// 	fmt.Printf("\tEnd: %s\n", window.EndTime)
	// }
}

func endVisibilityWindow(lastVisibilityWindow **VisibilityWindow, visibilityWindows *[]VisibilityWindow, endTime time.Time) {
	(**lastVisibilityWindow).EndTime = endTime
	*visibilityWindows = append(*visibilityWindows, **lastVisibilityWindow)
	*lastVisibilityWindow = nil
}

func isVisible(objectAltitute float64, objectAzimuth float64, config *Config) bool {
	alphaMin, alphaMax := getTelescopeMinMaxAltitute(config, objectAzimuth)
	return objectAltitute > Rad2deg(alphaMin) && objectAltitute < Rad2deg(alphaMax)
}

func getTelescopeMinMaxAltitute(config *Config, objectAzimuth float64) (float64, float64) {
	angleDiff := Deg2rad(math.Abs(objectAzimuth - config.DirectAzimuth))
	alphaMin := altitudeAtAzimuthDiff(config.FenceHeight-config.TelescopeHeight, config.DistanceToFence, angleDiff) - Deg2rad(3.0)
	alphaMax := altitudeAtAzimuthDiff(config.WindowHeight+config.FenceHeight-config.TelescopeHeight, config.DistanceToFence, angleDiff)
	return alphaMin, alphaMax
}

func altitudeAtAzimuthDiff(actualFenceHeight, distanceToFence, angleDiff float64) float64 {
	return math.Atan(actualFenceHeight * math.Cos(angleDiff) / distanceToFence)
}
