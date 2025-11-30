package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tps193/balcony-stargazer/internal/database"
	"github.com/tps193/balcony-stargazer/internal/visibility"
)

func main() {
	// logFile := initLogging()
	// defer logFile.Close()

	if len(os.Args) < 2 {
		fmt.Println("Expected 'observe' or 'suggest' subcommand")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "observe":
		runObserve(os.Args[2:])
	case "suggest":
		runSuggest(os.Args[2:])
	default:
		fmt.Println("expected 'observe' or 'suggest' subcommands")
		os.Exit(1)
	}

}

func runSuggest(s []string) {
	suggestCmd := flag.NewFlagSet("suggest", flag.ExitOnError)
	observationType := suggestCmd.String("observationtype", "", "Type of observation to suggest (e.g., planetary, deep-sky)")
	configFile := suggestCmd.String("configfile", "", "Path to the configuration file")
	configStr := suggestCmd.String("configstr", "", "String with configurations in JSON format")

	minVisibilityMin := suggestCmd.Int("minvisibilitytime", 0, "Minimum visibility duration in minutes")

	minSize := suggestCmd.Float64("minsize", -1.0, "Minimum size in arc minutes")
	maxSize := suggestCmd.Float64("maxsize", -1.0, "Maximum size in arc minutes")

	minMagnitude := suggestCmd.Float64("minmagnitude", -1.0, "Minimum magnitude")
	maxMagnitude := suggestCmd.Float64("maxmagnitude", -1.0, "Maximum magnitude")

	//TODO: make proper descriptions and add help
	timeFile := suggestCmd.String("timefile", "", "Path to the time file in RFC3339 format (e.g., 2024-06-30T22:30:00Z)")
	timeString := suggestCmd.String("timestr", "", "String with observation time windows in RFC3339 format (e.g., 2025-07-01T05:30:00Z)")

	logfile := suggestCmd.String("logfile", "", "Path to the log file")

	suggestCmd.Parse(s)

	f := initLogging(logfile)
	if f != nil {
		defer f.Close()
	}

	timeRanges, err := parseTime(timeFile, timeString)
	if err != nil {
		fmt.Println("Error parsing time range:", err)
		return
	}

	config, err := parseConfig(configFile, configStr)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	filter := database.Filter{
		MinSizeArcMinutes: *minSize,
		MaxSizeArcMinutes: *maxSize,
		MinMagnitude:      *minMagnitude,
		MaxMagnitude:      *maxMagnitude,
		ObjectType:        observationType,
	}

	//copilot put paths to catalog files
	catalogObjects, err := database.ParseCatalogCSV(filter, "/Users/sergey/Programming/GoProjects/balconyStargazer/database/NGC_with_common_names.csv")
	if err != nil {
		fmt.Println("Error parsing catalog CSV:", err)
		return
	}

	astroObjects, err := visibility.ToAstroObjects(catalogObjects)
	if err != nil {
		fmt.Println("Error converting catalog to astro objects:", err)
		return
	}

	visibilityInfos := visibility.CalculateAltitudeVisibility(astroObjects, config, timeRanges, 5, visibility.Filter{MinVisibilityDurationMinutes: *minVisibilityMin}, true)
	fmt.Println(visibility.NewSimpleOutputResult().Get(&visibilityInfos))
}

func runObserve(s []string) {
	observeCmd := flag.NewFlagSet("observe", flag.ExitOnError)
	configFile := observeCmd.String("configfile", "", "Path to the configuration file")
	configStr := observeCmd.String("configstr", "", "String with configurations in JSON format")

	objectFile := observeCmd.String("objectfile", "", "Path to the object file")
	objectStr := observeCmd.String("objectstr", "", "String with objects in JSON format")

	//TODO: make proper descriptions and add help
	timeFile := observeCmd.String("timefile", "", "Path to the time file in RFC3339 format (e.g., 2024-06-30T22:30:00Z)")
	timeString := observeCmd.String("timestr", "", "String with observation time windows in RFC3339 format (e.g., 2025-07-01T05:30:00Z)")

	minVisibilityMin := observeCmd.Int("minvisibilitytime", 0, "Minimum visibility duration in minutes")

	logfile := observeCmd.String("logfile", "", "Path to the log file")

	observeCmd.Parse(s)

	f := initLogging(logfile)
	if f != nil {
		defer f.Close()
	}

	config, err := parseConfig(configFile, configStr)
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	astroObjectValue, err := readFlag(objectFile, objectStr, "astronomical object")
	if err != nil {
		fmt.Println("Error reading astronomical object:", err)
		return
	}
	var objectsArray visibility.AstroObjectArray
	err = json.Unmarshal([]byte(astroObjectValue), &objectsArray)
	if err != nil {
		fmt.Println("Error parsing json:", err)
		return
	}

	timeRanges, err := parseTime(timeFile, timeString)
	if err != nil {
		fmt.Println("Error parsing time range:", err)
		return
	}

	log.Println(config)
	log.Println(objectsArray)

	visibilityInfos := visibility.CalculateAltitudeVisibility(&objectsArray, config, timeRanges, 5, visibility.Filter{MinVisibilityDurationMinutes: *minVisibilityMin}, true)
	fmt.Println(visibility.NewSimpleOutputResult().Get(&visibilityInfos))
}

func parseConfig(configFile, configStr *string) (*visibility.ConfigArray, error) {
	configValue, err := readFlag(configFile, configStr, "config")
	if err != nil {
		return nil, fmt.Errorf("Error reading configuration: %w", err)
	}
	var config visibility.ConfigArray
	err = json.Unmarshal([]byte(configValue), &config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing configuration: %w", err)
	}
	return &config, nil
}

func parseTime(timeFile, timeStr *string) ([]visibility.TimeRange, error) {
	if timeFile == nil && timeStr == nil {
		return nil, errors.New("time value is required")
	}

	type TimeRangeStr struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	}
	var timeRangesJson []TimeRangeStr
	if timeFile != nil && *timeFile != "" {
		fileContent, err := os.ReadFile(*timeFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading time file: %w", err)
		}
		err = json.Unmarshal(fileContent, &timeRangesJson)
		if err != nil {
			return nil, fmt.Errorf("Error parsing time file: %w", err)
		}
	} else if timeStr != nil && *timeStr != "" {
		err := json.Unmarshal([]byte(*timeStr), &timeRangesJson)
		if err != nil {
			return nil, fmt.Errorf("Error parsing time string: %w", err)
		}
	}

	var timeRanges []visibility.TimeRange
	if len(timeRangesJson) == 0 {
		return nil, errors.New("no valid time range found")
	} else {
		for _, tr := range timeRangesJson {
			startTime, err := time.Parse(time.RFC3339, tr.StartTime)
			if err != nil {
				return nil, fmt.Errorf("error parsing start time: %w", err)
			}
			endTime, err := time.Parse(time.RFC3339, tr.EndTime)
			if err != nil {
				return nil, fmt.Errorf("error parsing end time: %w", err)
			}
			if endTime.Before(startTime) {
				return nil, fmt.Errorf("end time %s is before start time %s", tr.EndTime, tr.StartTime)
			}
			timeRanges = append(timeRanges, visibility.TimeRange{StartTime: startTime, EndTime: endTime})
		}
	}
	return timeRanges, nil
}

func readFlag(fileValue, strValue *string, name string) (string, error) {
	if *fileValue == "" && *strValue == "" {
		return "", errors.New("no " + name + " provided")
	}
	if *fileValue != "" && *strValue != "" {
		return "", errors.New("either a file or a string can be provided for " + name)
	}

	result := *strValue
	if *fileValue != "" {
		file, err := os.ReadFile(*fileValue)
		if err != nil {
			return "", errors.New("Error reading " + name + " file: " + err.Error())
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
