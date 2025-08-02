# Balcony Stargazer

This tool is designed for astrophotographers with limited observation windows due to balconies or windows. It calculates object visibility based on your location, telescope position, and any obstructions like walls or fences, helping you determine the optimal time to observe. It also has MCP server implementation so that can be connected to LLM models.

# Example

## Command line

These files provide the object details and configuration for the observation setup.

**object.json:** Defines the astronomical object you want to observe, including its name, right ascension (RA), and declination (Dec).
```json
{
  "name": "Ghost of Cassiopeia",
  "ra": {
    "hour": 0,
    "min": 58,
    "sec": 26

  },
  "dec": {
    "degree": 60,
    "min": 53,
    "sec": 0
  }
}
```

**config.json:** Specifies the parameters of your observation location and telescope setup, such as fence height, window height, distance to the fence, telescope height, direct azimuth, and geographical position (latitude and longitude). It also defines the left and right azimuth limits for observations.
```json
{
  "fenceHeight": 43.25,
  "windowHeight": 62.0,
  "distanceToFence": 35,
  "telescopeHeight": 18.0,
  "directAzimuth": 80.0,
  "position": {
    "latitude": 37.38,
    "longitude": -121.89
  },
    "leftAzimuthLimit": 13.0,
    "rightAzimuthLimit": 120.0
}
```

```
$ balconystargazer -configfile=./config.json -objectfile=object.json -starttime=2025-07-30T22:13:00-07:00 -endtime=2025-07-31T05:23:00-07:00

$ Visibility of Ghost of Cassiopeia:
0: 4h10m0s
        Start: 2025-07-30 22:13:00 -0700 PDT (29.546459°)
        End: 2025-07-31 02:23:00 -0700 PDT (58.779733°)
```

## LLM

``` 
>> Calculate visibility of Ghost of Cassiopeia today from 10:30pm to 5:30am.

<< Visibility of Ghost of Cassiopeia:
0: 4h10m0s
        Start: 2025-07-30 22:13:00 -0700 PDT (29.546459°)
        End: 2025-07-31 02:23:00 -0700 PDT (58.779733°)
```

# Requirements

**General:**
1. Ruler (to measure your balcony parameters)
1. Golang 1.24.2

**For MCP server:**
1. AI client working with MCP Tools (for example [mcphost](https://github.com/mark3labs/mcphost))
1. LLM working with MCP tools (in case of local installation i.e. ollama:qwen2.5)

Note: it is better to use more advanced model as ollama:qwen2.5 doesn't know correct astronomical objects coordinates.

# Build

## Commmand line tool

```
go build -o balconystargazer ./cmd/main
```

## MCP Server

```
go build -o mcp ./cmd/mcp
```

# Run

## Command line tool

The command line tool requires following data to be provided:
1. Configuration json
1. Object information json
1. Start observation time with local timezone
1. End observation time with local timezone

### Configuration 

Configuration can be either taken from file (`-configfile`) or provided as a string literal (`-configstr`)

#### Structure

```json
{
  "fenceHeight": "number",
  "windowHeight": "number",
  "distanceToFence": "number",
  "telescopeHeight": "number",
  "directAzimuth": "degree",
  "position": {
    "latitude": "degree",
    "longitude": "degree"
  },
    "leftAzimuthLimit": "degree",
    "rightAzimuthLimit": "degree"
}
```

| Name | Data Type | Description |
|------|-----------|-------------|
| fenceHeight | number | Height of the fence or obstruction in front of the observation point |
| windowHeight | number | Height of the window from the top of the fence to the ceiling |
| distanceToFence | number | Distance from the observation point to the fence or obstruction |
| telescopeHeight | number | Height of the telescope above the ground level (from ground to the rotation point, for example to the servoe of Vaonis Vespera) |
| directAzimuth | number (degree) | Direct azimuth angle of the observation direction (traversal to the fence) |
| position.latitude | number (degree) | Geographic latitude of the observation location |
| position.longitude | number (degree) | Geographic longitude of the observation location |
| leftAzimuthLimit | number (degree) | Left boundary azimuth limit for observations |
| rightAzimuthLimit | number (degree) | Right boundary

#### Examples

```bash
balconystargazer -configfile=config.json -objectfile=object.json -starttime=2025-07-30T22:13:00-07:00 -endtime=2025-07-31T05:23:00-07:00
```

```bash
balconystargazer -configstr='{"fenceHeight": 43.25, "windowHeight": 62.0, "distanceToFence": 35, "telescopeHeight": 18.0, "directAzimuth": 80.0, "position": {"latitude": 37.38, "longitude": -121.89}, "leftAzimuthLimit": 13.0, "rightAzimuthLimit": 120.0}' -objectfile=object.json -starttime=2025-07-30T22:13:00-07:00 -endtime=2025-07-31T05:23:00-07:00
```

### Object Information

Object information can be either taken from file (`-objectfile`) or provided as a string literal (`-objectstr`)

#### Structure

```json
{
  "name": "string",
  "ra": {
    "hour": "number",
    "min": "number",
    "sec": "number"
  },
  "dec": {
    "degree": "number",
    "min": "number",
    "sec": "number"
  }
}
```

| Name | Data Type | Description |
|------|-----------|-------------|
| name | string | Name of the astronomical object |
| ra.hour | number | Hour component of right ascension (0-23) |
| ra.min | number | Minute component of right ascension (0-59) |
| ra.sec | number | Second component of right ascension (0-59) |
| dec.degree | number | Degree component of declination (-90 to +90) |
| dec.min | number | Minute component of declination (0-59) |
| dec.sec | number | Second component of declination (0-59) |

#### Examples

```bash
balconystargazer -configfile=config.json -objectfile=object.json -starttime=2025-07-30T22:13:00-07:00 -endtime=2025-07-31T05:23:00-07:00
```

```bash
balconystargazer -configfile=config.json -objectstr='{"name": "Ghost of Cassiopeia", "ra": {"hour": 0, "min": 58, "sec": 26}, "dec": {"degree": 60, "min": 53, "sec": 0}}' -starttime=2025-07-30T22:13:00-07:00 -endtime=2025-07-31T05:23:00-07:00
```

### Time

Start time and end time specify the window in which the tool would check for the object visibility.
The proper timezone must be specified in the value. Otherwise the calculation of local azimuth and altitude of the astronomical object at the specific time will be incorrect.

#### Examples

**Pacific Time (Los Angeles, California):**
```bash
-starttime=2025-07-30T22:30:00-07:00 -endtime=2025-07-31T05:30:00-07:00
```

**Eastern Time (New York, New York):**
```bash
-starttime=2025-07-30T22:30:00-04:00 -endtime=2025-07-31T05:30:00-04:00
```

**Central European Time (Berlin, Germany):**
```bash
-starttime=2025-07-30T22:30:00+02:00 -endtime=2025-07-31T05:30:00+02:00
```

**Japan Standard Time (Tokyo, Japan):**
```bash
-starttime=2025-07-30T22:30:00+09:00 -endtime=2025-07-31T05:30:00+09:00
```

### Logging

You can specify log file location using `-logfile` parameter

```bash
-logfile=test.log
```

```bash
balconystargazer -configfile=config.json -objectfile=object.json -starttime=2025-07-30T22:30:00-07:00 -endtime=2025-07-31T05:30:00-07:00 -logfile=stargazer.log
```

## MCP Server

The application is made with [mcp-go](https://github.com/mark3labs/mcp-go) project. The example of running MCP Server integration will be based on [mcphost](https://github.com/mark3labs/mcphost) client.

The application should be compiled as
```bash
go build -o mcp ./cmd/mcp
```

To configure the MCP server, you'll need to set up a configuration file for `mcphost` that points to the `mcp` executable.

The tool should be preconfigured in `mcphost` yaml congifuration

**mcpconfig.yaml**
```yaml
mcpServers:
  balconyStargazer:
    command: /path/to/executable/mcp
    name: balconyStargazer
```

Run client:
```bash
mcphost --model ollama:qwen2.5:latest --config mcpconfig.yaml
```

**Example request for local LLM**
```
calculate visibility of {
  "name": "Ghost of Cassiopeia",
  "ra": {
    "hour": 0,
    "min": 58,
    "sec": 26

  },
  "dec": {
    "degree": 60,
    "min": 53,
    "sec": 0
  }
}
using {
  "fenceHeight": 43.25,
  "windowHeight": 62.0,
  "distanceToFence": 35,
  "telescopeHeight": 18.0,
  "directAzimuth": 80.0,
  "position": {
    "latitude": 37.38,
    "longitude": -121.89
  },
    "leftAzimuthLimit": 13.0,
    "rightAzimuthLimit": 120.0
}
from 07-30-2025 22:30 to 07-31-2025 5:00. Set timezone based on the coordinates in the config.
```


That's it! When asking for an advice about object visibility the model should make a request to `balconyStargazer` app.