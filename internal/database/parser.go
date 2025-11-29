package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type CatalogRow struct {
	Name        string
	Type        string
	RA          string
	Dec         string
	Const       string
	MajAx       string
	MinAx       string
	PosAng      string
	BMag        string
	VMag        string
	Commonnames string
	// ... add more fields as needed
}

// ParseCatalogCSV parses a catalog CSV file and returns a slice of CatalogRow
func ParseCatalogCSV(filter Filter, filePath string) ([]CatalogRow, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1 // allow variable number of fields

	// Read header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var entries []CatalogRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		entry := CatalogRow{
			Name:        getField(record, 0),
			Type:        getField(record, 1),
			RA:          getField(record, 2),
			Dec:         getField(record, 3),
			Const:       getField(record, 4),
			MajAx:       getField(record, 5),
			MinAx:       getField(record, 6),
			PosAng:      getField(record, 7),
			BMag:        getField(record, 8),
			VMag:        getField(record, 9),
			Commonnames: getField(record, 28),
		}

		log.Println("Processing entry:", entry.Name)

		objectType := entry.Type
		if filter.ObjectType != nil && *filter.ObjectType != "" && objectType != *filter.ObjectType {
			continue
		}

		if filter.MinSizeArcMinutes > 0 || filter.MaxSizeArcMinutes > 0 {
			majAx, err := toFloat(entry.MajAx)
			if err != nil {
				continue
			}
			minAx, err := toFloat(entry.MinAx)
			if err != nil {
				continue
			}
			minSize := min(majAx, minAx)
			maxSize := max(majAx, minAx)

			if filter.MinSizeArcMinutes > 0 && minSize < filter.MinSizeArcMinutes {
				continue
			}
			if filter.MaxSizeArcMinutes > 0 && maxSize > filter.MaxSizeArcMinutes {
				continue
			}
		}

		if filter.MinMagnitude > 0 || filter.MaxMagnitude > 0 {
			vMag, err := toFloat(entry.VMag)
			if err != nil {
				continue
			}

			if filter.MinMagnitude > 0 && vMag > filter.MinMagnitude {
				continue
			}
			if filter.MaxMagnitude > 0 && vMag < filter.MaxMagnitude {
				continue
			}
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

func toFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

// getField safely gets a field from a record or returns an empty string if out of range
func getField(record []string, idx int) string {
	if idx < len(record) {
		return record[idx]
	}
	return ""
}
