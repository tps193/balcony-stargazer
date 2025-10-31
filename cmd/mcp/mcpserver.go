package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tps193/balcony-stargazer/internal/visibility"
)

const (
	AstroObjects = "astroObjects"
	Config       = "config"
)

func main() {
	f, err := os.OpenFile("mcp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Set log output to the file
	log.SetOutput(f)

	log.Println("Log initialized")

	defer f.Close()

	// Create a new MCP server
	s := server.NewMCPServer(
		"Balcony Stargzer",
		"0.0.1",
		server.WithToolCapabilities(true),
	)
	schema := jsonschema.Reflect(&visibility.AstroObjectArray{})
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		fmt.Println("Error creating schema for AstroObjectArray:", err)
		panic(err)
	}
	astroObjectSchema := string(schemaBytes)

	schema = jsonschema.Reflect(&visibility.ConfigArray{})
	schemaBytes, err = json.Marshal(schema)
	if err != nil {
		fmt.Println("Error creating schema for ConfigArray:", err)
		panic(err)
	}
	configSchema := string(schemaBytes)

	// Add visibilityInWindowTool
	visibilityInWindowTool := mcp.NewTool("astro_object_visibility",
		mcp.WithDescription("Allows to calculate visibility windows for astronomical object within specified range. Ask user for parameters and wait input before running the tool."),
		mcp.WithString(AstroObjects,
			mcp.Required(),
			mcp.Description("Name and coordinates of the array of astronomical object formatted as single string json "+astroObjectSchema),
		),
		mcp.WithString(Config,
			mcp.Required(),
			mcp.Description("Must be asked from user. Configuration for visibility calculation formatted as single string json "+configSchema),
		),
		mcp.WithString("startTime",
			mcp.Required(),
			mcp.Description("Must be asked from user and not generated. Observation start time in RFC3339 format (e.g., 2024-06-30T22:30:00-05:00). Timezone is required and must be calculated from the user location from config parameter."),
		),
		mcp.WithString("endTime",
			mcp.Required(),
			mcp.Description("Must be asked from user and not generated. Observation end time in RFC3339 format (e.g., 2025-07-01T05:30:00-05:00). Timezone is required and must be calculated from the user location from config parameter."),
		),
	)

	quickVisibilityFilterTool := mcp.NewTool("quick_visibility_filter",
		mcp.WithDescription("Quickly checks if the object is ever visible from the given location and if it ever comes into the specified azimuth window. Does not require time range. Should be used especially when user asks not for a specific object visibility, but to find an object that is visible from the given location and within the azimuth window. For example: find me an object that is visible from my location and within azimuth window 90°-270°. Ask user for parameters and wait input before running the tool."),
		mcp.WithString(AstroObjects,
			mcp.Required(),
			mcp.Description("Name and coordinates of the astronomical object formatted as single string json "+astroObjectSchema),
		),
		mcp.WithString(Config,
			mcp.Required(),
			mcp.Description("Must be asked from user. Configuration for visibility calculation formatted as single string json "+configSchema),
		),
	)

	// Add tool handler
	s.AddTool(visibilityInWindowTool, visibilityHandler)
	s.AddTool(quickVisibilityFilterTool, quickVisibilityFilterHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func visibilityHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jsonStr, err := request.RequireString(AstroObjects)
	if err != nil {
		log.Println("Error requiring astroObjectInfo: ", err)
		return mcp.NewToolResultError("Error: error getting required parameter " + AstroObjects + " due to " + err.Error() + ". Single line json string is expected"), nil
	}
	log.Println("AstroObjectInfo: ", jsonStr)

	astroObjectArray := &visibility.AstroObjectArray{}
	err = json.Unmarshal([]byte(jsonStr), astroObjectArray)
	if err != nil {
		log.Println(err.Error())
		return mcp.NewToolResultError("Error unmarshalling Astro Object json: " + err.Error()), nil
	}
	log.Println("Unmarshalled object: ", astroObjectArray)

	jsonStr, err = request.RequireString(Config)
	if err != nil {
		log.Println("Error requiring config:", err)
		return mcp.NewToolResultError(err.Error() + ". Single line string json is expected"), nil
	}
	log.Println("Config: ", jsonStr)

	config := &visibility.ConfigArray{}
	err = json.Unmarshal([]byte(jsonStr), config)
	if err != nil {
		log.Println(err.Error())
		return mcp.NewToolResultError(err.Error()), nil
	}
	log.Println("Unmarshalled config: ", config)

	startTimeStr, err := request.RequireString("startTime")
	if err != nil {
		log.Println("Error requiring startTime:", err)
		return mcp.NewToolResultError(err.Error() + ". Start time in RFC3339 format (e.g., 2024-06-30T22:30:00Z) is expected"), nil
	}
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		log.Println("Error parsing start time:", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	endTimeStr, err := request.RequireString("endTime")
	if err != nil {
		log.Println("Error requiring endTime:", err)
		return mcp.NewToolResultError(err.Error() + ". End time in RFC3339 format (e.g., 2025-07-01T05:30:00Z) is expected"), nil
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		log.Println("Error parsing end time:", err)
		return mcp.NewToolResultError(err.Error()), nil
	}
	log.Printf("Observed time from %s to %s\n", startTime, endTime)

	visibilityInfos := visibility.CalculateAltitudeVisibility(astroObjectArray, config, startTime, endTime, 5, true)
	result := visibility.NewJsonOutput().Get(visibilityInfos)
	return mcp.NewToolResultText(result), nil
}

func quickVisibilityFilterHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jsonStr, err := request.RequireString(AstroObjects)
	if err != nil {
		log.Println("Error requiring astroObjectInfo: ", err)
		return mcp.NewToolResultError("Error: error getting required parameter " + AstroObjects + " due to " + err.Error() + ". Single line json string is expected"), nil
	}
	log.Println("AstroObjectInfo: ", jsonStr)

	astroObject := visibility.AstroObject{}
	err = json.Unmarshal([]byte(jsonStr), &astroObject)
	if err != nil {
		log.Println(err.Error())
		return mcp.NewToolResultError("Error unmarshalling Astro Object json: " + err.Error()), nil
	}
	log.Println("Unmarshalled object: ", astroObject)

	jsonStr, err = request.RequireString(Config)
	if err != nil {
		log.Println("Error requiring config:", err)
		return mcp.NewToolResultError(err.Error() + ". Single line string json is expected"), nil
	}
	log.Println("Config: ", jsonStr)

	config := &visibility.Config{}
	err = json.Unmarshal([]byte(jsonStr), config)
	if err != nil {
		log.Println(err.Error())
		return mcp.NewToolResultError(err.Error()), nil
	}
	log.Println("Unmarshalled config: ", config)

	neverVisible := visibility.ObjectNeverVisible(astroObject, config)
	everInAzimuthWindow := visibility.ObjectEverInAzimuthWindow(astroObject, config)

	result := fmt.Sprintf("Object %s never visible: %t, ever in azimuth window: %t", astroObject.Name, neverVisible, everInAzimuthWindow)
	return mcp.NewToolResultText(result), nil
}
