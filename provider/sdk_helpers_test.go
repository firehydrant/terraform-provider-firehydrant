package provider

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalLabels(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    map[string]string
		expectError bool
	}{
		{
			name:     "nil labels",
			input:    nil,
			expected: nil,
		},
		{
			name: "string values",
			input: func() interface{} {
				// Simulate SDK struct with string labels
				var data struct {
					Key1 string `json:"key1"`
					Key2 string `json:"key2"`
				}
				data.Key1 = "value1"
				data.Key2 = "value2"
				return data
			}(),
			expected: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "numeric values converted to strings",
			input: func() interface{} {
				// Simulate SDK struct with numeric labels
				jsonStr := `{"count": 123, "price": 45.67}`
				var result map[string]interface{}
				json.Unmarshal([]byte(jsonStr), &result)
				return result
			}(),
			expected: map[string]string{
				"count": "123",
				"price": "45.67",
			},
		},
		{
			name: "empty struct results in nil",
			input: struct {
			}{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unmarshalLabels(tt.input)
			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expected == nil && result != nil {
				t.Errorf("expected nil, got %v", result)
			}
			if tt.expected != nil && result == nil {
				t.Errorf("expected %v, got nil", tt.expected)
			}
			if tt.expected != nil && result != nil {
				if len(result) != len(tt.expected) {
					t.Errorf("expected %d labels, got %d", len(tt.expected), len(result))
				}
				for key, expectedValue := range tt.expected {
					if actualValue, ok := result[key]; !ok {
						t.Errorf("missing key %s", key)
					} else if actualValue != expectedValue {
						t.Errorf("for key %s: expected %s, got %s", key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}
