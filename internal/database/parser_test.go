package database

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestParseCatalogCSV_PrintJSON(t *testing.T) {
	files := []string{"../../database/NGC_with_common_names.csv", "../../database/addendum.csv"}
	for _, file := range files {
		t.Logf("Parsing %s", file)
		start := time.Now()
		entries, err := ParseCatalogCSV(Filter{}, file)
		duration := time.Since(start)
		t.Logf("Parsing %s took %v", file, duration)
		if err != nil {
			t.Fatalf("Failed to parse %s: %v", file, err)
		}
		for i, entry := range entries {
			jsonLine, err := json.Marshal(entry)
			if err != nil {
				t.Errorf("Failed to marshal entry %d in %s: %v", i, file, err)
				continue
			}
			fmt.Println(string(jsonLine))
		}
	}
}

func TestParseCatalogCSV_Filters(t *testing.T) {
	file := "../../database/NGC_with_common_names.csv"
	tests := []struct {
		name    string
		filter  Filter
		wantAny bool
	}{
		{
			name:    "No filter",
			filter:  Filter{},
			wantAny: true,
		},
		{
			name:    "ObjectType=G",
			filter:  Filter{ObjectType: strPtr("G")},
			wantAny: true,
		},
		{
			name:    "MinMagnitude=10",
			filter:  Filter{MinMagnitude: 10.0},
			wantAny: true,
		},
		{
			name:    "MaxMagnitude=8",
			filter:  Filter{MaxMagnitude: 8.0},
			wantAny: true,
		},
		{
			name:    "MinSizeArcMinutes=1",
			filter:  Filter{MinSizeArcMinutes: 1.0},
			wantAny: true,
		},
		{
			name:    "MaxSizeArcMinutes=0.5",
			filter:  Filter{MaxSizeArcMinutes: 0.5},
			wantAny: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := ParseCatalogCSV(tt.filter, file)
			if err != nil {
				t.Fatalf("ParseCatalogCSV failed: %v", err)
			}
			if tt.wantAny && len(entries) == 0 {
				t.Errorf("Expected at least one entry for filter %+v, got 0", tt.filter)
			}
			t.Logf("Filter: %+v, got %d entries", tt.filter, len(entries))
		})
	}
}

func strPtr(s string) *string { return &s }
