package database

import (
	"encoding/csv"
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
func ParseCatalogCSV(objectType, filePath string) ([]CatalogRow, error) {
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
		if err != nil {
			break // EOF
		}
		entry := CatalogRow{
			Name:   getField(record, 0),
			Type:   getField(record, 1),
			RA:     getField(record, 2),
			Dec:    getField(record, 3),
			Const:  getField(record, 4),
			MajAx:  getField(record, 5),
			MinAx:  getField(record, 6),
			PosAng: getField(record, 7),
			BMag:   getField(record, 8),
			VMag:   getField(record, 9),
			//copilot: check correct column number for Commonnames: Name;Type;RA;Dec;Const;MajAx;MinAx;PosAng;B-Mag;V-Mag;J-Mag;H-Mag;K-Mag;SurfBr;Hubble;Pax;Pm-RA;Pm-Dec;RadVel;Redshift;Cstar U-Mag;Cstar B-Mag;Cstar V-Mag;M;NGC;IC;Cstar Names;Identifiers;Commonnames;NED notes;OpenNGC notes;Sources
			Commonnames: getField(record, 29),
			// ... add more fields as needed
		}
		if objectType == "" || entry.Type == objectType {
			entries = append(entries, entry)
		}
	}
	return entries, nil
}

// getField safely gets a field from a record or returns an empty string if out of range
func getField(record []string, idx int) string {
	if idx < len(record) {
		return record[idx]
	}
	return ""
}
