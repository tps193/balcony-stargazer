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
		entries, err := ParseCatalogCSV("", file)
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
