package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"balcony-stargazer/internal/visibility"
)

func main() {
	// logFile := initLogging()
	// defer logFile.Close()

	configFile := flag.String("configfile", "", "Path to the configuration file")
	configStr := flag.String("configstr", "", "Configuration string in JSON format")

	objectFile := flag.String("objectfile", "", "Path to the object file")
	objectStr := flag.String("objectstr", "", "Object string in JSON format")

	startTimeValue := flag.String("starttime", "", "Start time in RFC3339 format (e.g., 2024-06-30T22:30:00Z)")
	endTimeValue := flag.String("endtime", "", "End time in RFC3339 format (e.g., 2025-07-01T05:30:00Z)")

	logfile := flag.String("logfile", "", "Path to the log file")

	flag.Parse()

	f := initLogging(logfile)
	if f != nil {
		defer f.Close()
	}

	configValue, err := readFlag(configFile, configStr, "config")
	if err != nil {
		fmt.Println("Error reading configuration:", err)
		return
	}
	var config visibility.Config
	err = json.Unmarshal([]byte(configValue), &config)
	if err != nil {
		fmt.Println("Error parsing configuration:", err)
		return
	}

	astroObjectValue, err := readFlag(objectFile, objectStr, "astronomical object")
	if err != nil {
		fmt.Println("Error reading astronomical object:", err)
		return
	}
	var object visibility.AstroObject
	err = json.Unmarshal([]byte(astroObjectValue), &object)
	if err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}

	startTime, err := parseTime(startTimeValue)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return
	}
	endTime, err := parseTime(endTimeValue)
	if err != nil {
		fmt.Println("Error parsing end time:", err)
		return
	}

	log.Printf("Observed time from %s to %s\n", startTime, endTime)
	log.Println(config)
	log.Println(object)
	log.Printf("Start time: %s, End time: %s\n", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	visibilityWindows := visibility.CalculateAltitudeVisibility(&object, &config, startTime, endTime, 5, true)
	fmt.Println(visibility.NewSimpleOutputResult().Get(&object, &visibilityWindows))
}

func parseTime(timeStr *string) (time.Time, error) {
	if timeStr == nil || *timeStr == "" {
		return time.Time{}, errors.New("time value is required")
	}
	timeObj, err := time.Parse(time.RFC3339, *timeStr)
	if err != nil {
		return time.Time{}, err
	}
	return timeObj, nil
}

func readFlag(fileValue, strValue *string, name string) (string, error) {
	if *fileValue == "" && *strValue == "" {
		return "", errors.New("no configuration provided")
	}
	if *fileValue != "" && *strValue != "" {
		return "", errors.New("either a configuration file or a configuration string can be provided")
	}

	result := *strValue
	if *fileValue != "" {
		file, err := os.ReadFile(*fileValue)
		if err != nil {
			return "", errors.New("Error reading configuration file: " + err.Error())
		}
		result = string(file)
	}
	return result, nil
}

func initLogging(logfile *string) *os.File {
	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening log file:", err)
			return nil
		}
		log.SetOutput(f)
		return f
	} else {
		return nil
	}
}
