package provider

import (
	"encoding/json"
	"fmt"
)

// unmarshalLabels takes an SDK labels struct (which is an empty struct used to represent
// dynamic key-value pairs) and unmarshals it into a map[string]string.
// The SDK uses empty structs for fields with undefined structure. We need to intercept
// the unmarshalling and convert any non-string values to strings for Terraform compatibility.
func unmarshalLabels(labelsStruct interface{}) (map[string]string, error) {
	if labelsStruct == nil {
		return nil, nil
	}

	// Marshal the struct back to JSON to access the underlying data
	jsonBytes, err := json.Marshal(labelsStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal labels: %w", err)
	}

	// Unmarshal into a generic map first
	var rawMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &rawMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
	}

	// If the map is empty, return nil (no labels set)
	if len(rawMap) == 0 {
		return nil, nil
	}

	// Convert all values to strings
	stringMap := make(map[string]string, len(rawMap))
	for key, value := range rawMap {
		switch v := value.(type) {
		case string:
			stringMap[key] = v
		case nil:
			stringMap[key] = ""
		default:
			// Convert any non-string values to their JSON representation
			jsonValue, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal label value for key %s: %w", key, err)
			}
			stringMap[key] = string(jsonValue)
		}
	}

	return stringMap, nil
}
