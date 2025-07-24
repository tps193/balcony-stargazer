package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/alecthomas/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type AstroObject struct {
	Name string         `json:"name"`
	Ra   RightAscention `json:"ra"`
	Dec  Declanation    `json:"dec"`
}

type Position struct {
	Latitude   float64
	Longtitude float64
}

type RightAscention struct {
	Hour float64 `json:"hour"`
	Min  float64 `json:"min"`
	Sec  float64 `json:"sec"`
}

type Declanation struct {
	Degree float64 `json:"degree"`
	Min    float64 `json:"min"`
	Sec    float64 `json:"sec"`
}

func main() {
	f, err := os.OpenFile("", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	schema := jsonschema.Reflect(&AstroObject{})
	shemaBytes, err := json.Marshal(schema)
	if err != nil {
		fmt.Println("Error creating schema for AstroObject:", err)
		panic(err)
	}

	// Add tool
	tool := mcp.NewTool("astro_object_visibility",
		mcp.WithDescription("Allows to calculate visibility windows for astronomical objects"),
		mcp.WithString("astroObjectInfo",
			mcp.Required(),
			mcp.Description("Name and coordinates of the astronomical object formatted as single string json "+string(shemaBytes)),
		),
	)

	// calcule visibility of Veil Nebula using the visibility tool
	log.Println(string(shemaBytes))

	// Add tool handler
	s.AddTool(tool, helloHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// name, err := request.RequireString("name")
	// if err != nil {
	// 	return mcp.NewToolResultError(err.Error()), nil
	// }

	// return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
	//show visibility of Veil Nebula using the visibility tool

	// if request.Params.Arguments == nil {
	// 	return mcp.NewToolResultError("No coordinates json provided"), nil
	// }

	info, err := request.RequireString("astroObjectInfo")
	if err != nil {
		return mcp.NewToolResultError(err.Error() + ". Single line string json is expected"), nil
	}
	log.Println("astroObjectInfo:", info)

	visibilityTool := "/Users/sergey/Programming/GoProjects/balconyStargazer/main"
	astroObject := &AstroObject{}
	// request.BindArguments(astroObject)
	log.Println(astroObject)
	err = json.Unmarshal([]byte(info), astroObject)
	if err != nil {
		log.Println(err.Error())
		return mcp.NewToolResultError(err.Error()), nil
	}

	encoded, err := json.Marshal(astroObject)
	if err != nil {
		log.Println(err.Error())
		return mcp.NewToolResultError(err.Error()), nil
	}
	execCmd := exec.Command(visibilityTool, string(encoded))
	execOutput, err := execCmd.Output()
	if err != nil {
		return mcp.NewToolResultText(err.Error()), nil
	}
	return mcp.NewToolResultText(string(execOutput)), nil
}
