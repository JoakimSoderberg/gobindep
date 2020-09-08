package module

import (
	"fmt"
	"strings"
)

// Module represents a single Go module.
//
// Depending on the source that this is parsed from, fields may be empty.
// All helper functions on Module work with zero values. See their associated
// documentation for more information on exact behavior.
type Module struct {
	Path    string  `json:"path"`              // Import path, such as "github.com/mitchellh/golicense"
	Version string  `json:"version"`           // Version like "v1.2.3"
	Hash    string  `json:"hash"`              // Hash such as "h1:abcd1234"
	Replace *Module `json:"replace,omitempty"` // If the module was replaced
}

// String returns a human readable string format.
func (m *Module) String() string {
	return fmt.Sprintf("%s (%s)", m.Path, m.Version)
}

// ParseExeData parses the raw dependency information from a compiled Go
// binary's readonly data section. Any unexpected values will return errors.
func ParseExeData(raw string) ([]Module, error) {
	var result []Module
	for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
		row := strings.Split(line, "\t")

		// Ignore non-dependency information, such as path/mod. The
		// "=>" syntax means it is a replacement.
		if row[0] != "dep" && row[0] != "=>" {
			continue
		}

		if len(row) == 3 {
			// A row with 3 can occur if there is no hash data for the
			// dependency.
			row = append(row, "")
		}

		if len(row) != 4 {
			return nil, fmt.Errorf(
				"Unexpected raw dependency format: %s", line)
		}

		next := Module{
			Path:    row[1],
			Version: row[2],
			Hash:    row[3],
			Replace: nil,
		}

		// If this is a replacement, then add it to the last result
		if row[0] == "=>" {
			prev := &result[len(result)-1]
			prev.Replace = &next
		} else {
			// Not a replacement so append it to the list
			result = append(result, next)
		}
	}

	return result, nil
}
