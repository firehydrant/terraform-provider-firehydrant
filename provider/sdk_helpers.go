package provider

import (
	"encoding/json"
	"fmt"
)

// unmarshalLabels takes an SDK labels map (map[string]any and converts it
// into a map[string]string for Terraform compatibility.
// We need to convert any non-string values to strings since Terraform's TypeMap expects strings.
func unmarshalLabels(labelsMap interface{}) (map[string]string, error) {
	if labelsMap == nil {
		return nil, nil
	}

	switch labels := labelsMap.(type) {
	case map[string]any:
		// If the map is empty, return nil (no labels set)
		if len(labels) == 0 {
			return nil, nil
		}

		// Convert all values to strings
		stringMap := make(map[string]string, len(labels))
		for key, value := range labels {
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

	default:
		// Fallback: try to marshal/unmarshal for older SDK versions with empty structs
		jsonBytes, err := json.Marshal(labelsMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal labels: %w", err)
		}

		var rawMap map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &rawMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal labels: %w", err)
		}

		if len(rawMap) == 0 {
			return nil, nil
		}

		stringMap := make(map[string]string, len(rawMap))
		for key, value := range rawMap {
			switch v := value.(type) {
			case string:
				stringMap[key] = v
			case nil:
				stringMap[key] = ""
			default:
				jsonValue, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal label value for key %s: %w", key, err)
				}
				stringMap[key] = string(jsonValue)
			}
		}
		return stringMap, nil
	}
}
