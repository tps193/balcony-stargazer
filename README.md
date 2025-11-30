# Balcony Stargazer

This tool is designed for astrophotographers with limited observation windows due to balconies or windows. It calculates object visibility based on your location, telescope position, and any obstructions like walls or fences, helping you determine the optimal time to observe. It also has MCP server implementation so that can be connected to LLM models.

The database with the astroobjects for the planned is picked up from [https://github.com/mattiaverga/OpenNGC](https://github.com/mattiaverga/OpenNGC) project.

# Example

## Command line

The application supports two subcommands:
- **`observe`**: Calculate visibility for specific astronomical objects you provide
- **`suggest`**: Search the catalog and suggest observable objects matching your criteria

### Observe Subcommand

Calculate visibility for specific astronomical objects.

**Usage:**
```bash
./main observe -configfile=./config.json -objectfile=object.json -timefile=time.json
```

**object.json:** Defines astronomical objects as an array. Each object includes name, right ascension (RA), and declination (Dec).
```json
{
  "objects": [
    {
      "name": "Ghost of Cassiopeia",
      "ra": {
        "hour": 0,
        "min": 59,
        "sec": 59
      },
      "dec": {
        "degree": 60,
        "min": 0,
        "sec": 0
      }
    }
  ]
}
```

**config.json:** Specifies observation location and telescope setup as an array of configurations (supports multiple observation positions/windows).
```json
{
  "configs": [
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
  ]
}
```

**time.json:** Time windows as an array with RFC3339 formatted timestamps including timezone.
```json
[
  {
    "startTime": "2025-07-30T22:13:00-07:00",
    "endTime": "2025-07-31T05:23:00-07:00"
  }
]
```

**Example output:**
```
Visibility of Ghost of Cassiopeia:
0: 4h10m0s
        Start: 2025-07-30 22:13:00 -0700 PDT (29.546459°)
        End: 2025-07-31 02:23:00 -0700 PDT (58.779733°)
```

### Suggest Subcommand

Search the catalog and suggest observable objects matching your criteria.

**Usage:**
```bash
./main suggest -configfile=./config.json -timefile=time.json -observationtype=HII -minsize=5.0 -maxmagnitude=10.0
```

**Available filters:**
- `-observationtype`: Object type (e.g., `HII`, `G`, `Neb`, `OCl`, `PN`)
- `-minsize`, `-maxsize`: Size constraints in arc minutes
- `-minmagnitude`, `-maxmagnitude`: Magnitude constraints
- `-minvisibilitytime`: Minimum visibility duration in minutes

**Example output:**
```
Visibility of NGC7635 (Bubble Nebula):
0: 3h45m0s
        Start: 2025-07-30 22:30:00 -0700 PDT (35.2°)
        End: 2025-07-31 02:15:00 -0700 PDT (62.8°)
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

## Command line tool

```bash
go build -o main ./cmd/main
```

## MCP Server

```bash
go build -o mcp ./cmd/mcp
```

# Run

## Command Line Tool

The tool supports two subcommands: `observe` and `suggest`.

### Observe Command

Calculate visibility for specific astronomical objects.

**Syntax:**
```bash
./main observe [flags]
```

**Required flags (one of each pair):**
- `-configfile=<path>` or `-configstr=<json>`
- `-objectfile=<path>` or `-objectstr=<json>`
- `-timefile=<path>` or `-timestr=<json>`

**Optional flags:**
- `-minvisibilitytime=<minutes>`: Minimum visibility duration (default: 0)
- `-logfile=<path>`: Log file location

**Examples:**
```bash
# Using files
./main observe -configfile=config.json -objectfile=objects.json -timefile=time.json

# Using string literals
./main observe -configstr='{"configs":[{...}]}' -objectstr='{"objects":[{...}]}' -timestr='[{"startTime":"...","endTime":"..."}]'

# With logging and minimum visibility
./main observe -configfile=config.json -objectfile=objects.json -timefile=time.json -minvisibilitytime=30 -logfile=output.log
```

### Suggest Command

Search catalog and suggest observable objects matching criteria.

**Syntax:**
```bash
./main suggest [flags]
```

**Required flags (one of each pair):**
- `-configfile=<path>` or `-configstr=<json>`
- `-timefile=<path>` or `-timestr=<json>`

**Optional filter flags:**
- `-observationtype=<type>`: Object type (e.g., HII, G, Neb, OCl, PN)
- `-minsize=<arcmin>`: Minimum object size in arc minutes (use -1 to ignore)
- `-maxsize=<arcmin>`: Maximum object size in arc minutes (use -1 to ignore)
- `-minmagnitude=<mag>`: Minimum magnitude (use -1 to ignore)
- `-maxmagnitude=<mag>`: Maximum magnitude (use -1 to ignore)
- `-minvisibilitytime=<minutes>`: Minimum visibility duration (default: 0)
- `-logfile=<path>`: Log file location

**Examples:**
```bash
# Find all HII regions
./main suggest -configfile=config.json -timefile=time.json -observationtype=HII

# Find bright, large galaxies with at least 1 hour visibility
./main suggest -configfile=config.json -timefile=time.json -observationtype=G -minsize=5.0 -maxmagnitude=10.0 -minvisibilitytime=60

# Find planetary nebulae visible for at least 30 minutes
./main suggest -configfile=config.json -timefile=time.json -observationtype=PN -minvisibilitytime=30 -logfile=suggest.log
```

### Configuration Format

Configuration can be provided via file (`-configfile`) or string literal (`-configstr`).

**Structure:**
```json
{
  "configs": [
    {
      "fenceHeight": <number>,
      "windowHeight": <number>,
      "distanceToFence": <number>,
      "telescopeHeight": <number>,
      "directAzimuth": <degrees>,
      "position": {
        "latitude": <degrees>,
        "longitude": <degrees>
      },
      "leftAzimuthLimit": <degrees>,
      "rightAzimuthLimit": <degrees>
    }
  ]
}
```

**Field descriptions:**

| Name | Type | Description |
|------|------|-------------|
| configs | array | Array of configuration objects (supports multiple observation positions) |
| fenceHeight | number | Height of fence/obstruction in front of observation point |
| windowHeight | number | Height of window from top of fence to ceiling |
| distanceToFence | number | Distance from observation point to fence/obstruction |
| telescopeHeight | number | Height of telescope above ground (to rotation point) |
| directAzimuth | number (degrees) | Direct azimuth angle of observation direction (perpendicular to fence) |
| position.latitude | number (degrees) | Geographic latitude of observation location |
| position.longitude | number (degrees) | Geographic longitude of observation location |
| leftAzimuthLimit | number (degrees) | Left boundary azimuth limit for observations |
| rightAzimuthLimit | number (degrees) | Right boundary azimuth limit for observations |

**Examples:**

Single configuration:
```json
{
  "configs": [
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
  ]
}
```

Multiple configurations (multiple obstructions):
```json
{
  "configs": [
    {
      "fenceHeight": 43.25,
      "windowHeight": 62.0,
      "distanceToFence": 12,
      "telescopeHeight": 45.0,
      "directAzimuth": 66.0,
      "position": {
        "latitude": 37.38,
        "longitude": -121.89
      },
      "leftAzimuthLimit": 350.0,
      "rightAzimuthLimit": 160.0
    },
    {
      "fenceHeight": 43.25,
      "windowHeight": 62.0,
      "distanceToFence": 33,
      "telescopeHeight": 45.0,
      "directAzimuth": 340.0,
      "position": {
        "latitude": 37.38,
        "longitude": -121.89
      },
      "leftAzimuthLimit": 255.0,
      "rightAzimuthLimit": 298.0
    }
  ]
}
```

### Object Information Format (observe command only)

Object information can be provided via file (`-objectfile`) or string literal (`-objectstr`).

**Structure:**
```json
{
  "objects": [
    {
      "name": <string>,
      "ra": {
        "hour": <number>,
        "min": <number>,
        "sec": <number>
      },
      "dec": {
        "degree": <number>,
        "min": <number>,
        "sec": <number>
      },
      "objectType": <string (optional)>
    }
  ]
}
```

**Field descriptions:**

| Name | Type | Description |
|------|------|-------------|
| objects | array | Array of astronomical objects to observe |
| name | string | Name of the astronomical object |
| ra.hour | number | Hour component of right ascension (0-23) |
| ra.min | number | Minute component of right ascension (0-59) |
| ra.sec | number | Second component of right ascension (0-59) |
| dec.degree | number | Degree component of declination (-90 to +90) |
| dec.min | number | Minute component of declination (0-59) |
| dec.sec | number | Second component of declination (0-59) |
| objectType | string (optional) | Type of object (HII, G, Neb, etc.) |

**Example:**
```json
{
  "objects": [
    {
      "name": "Andromeda Galaxy (M31)",
      "ra": {
        "hour": 0,
        "min": 42,
        "sec": 44
      },
      "dec": {
        "degree": 41,
        "min": 16,
        "sec": 9
      },
      "objectType": "G"
    },
    {
      "name": "Orion Nebula (M42)",
      "ra": {
        "hour": 5,
        "min": 35,
        "sec": 17
      },
      "dec": {
        "degree": -5,
        "min": 23,
        "sec": 28
      },
      "objectType": "HII"
    }
  ]
}
```

### Time Format

Time windows can be provided via file (`-timefile`) or string literal (`-timestr`).

**Structure:**
```json
[
  {
    "startTime": <RFC3339 timestamp with timezone>,
    "endTime": <RFC3339 timestamp with timezone>
  }
]
```

The format must be RFC3339 with timezone offset (e.g., `-07:00` for PDT, `+02:00` for CEST).

**Examples:**

**Pacific Time (Los Angeles, California):**
```json
[
  {
    "startTime": "2025-07-30T22:30:00-07:00",
    "endTime": "2025-07-31T05:30:00-07:00"
  }
]
```

**Eastern Time (New York, New York):**
```json
[
  {
    "startTime": "2025-07-30T22:30:00-04:00",
    "endTime": "2025-07-31T05:30:00-04:00"
  }
]
```

**Central European Time (Berlin, Germany):**
```json
[
  {
    "startTime": "2025-07-30T22:30:00+02:00",
    "endTime": "2025-07-31T05:30:00+02:00"
  }
]
```

**Japan Standard Time (Tokyo, Japan):**
```json
[
  {
    "startTime": "2025-07-30T22:30:00+09:00",
    "endTime": "2025-07-31T05:30:00+09:00"
  }
]
```

**Multiple time windows:**
```json
[
  {
    "startTime": "2025-07-30T22:00:00-07:00",
    "endTime": "2025-07-31T02:00:00-07:00"
  },
  {
    "startTime": "2025-07-31T22:00:00-07:00",
    "endTime": "2025-08-01T02:00:00-07:00"
  }
]
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