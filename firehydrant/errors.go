package firehydrant

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorNotFound     = errors.New("resource not found")
	ErrorUnauthorized = errors.New("unauthorized")
)

type APIError struct {
	Error    string   `json:"error"`
	Detail   string   `json:"detail"`
	Messages []string `json:"messages"`
}

// UnmarshalJSON custom unmarshaler to handle messages field that can be either a string or array of strings
func (e *APIError) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to handle both formats
	var temp struct {
		Error    string      `json:"error"`
		Detail   string      `json:"detail"`
		Messages interface{} `json:"messages"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	e.Error = temp.Error
	e.Detail = temp.Detail

	// Handle messages field - can be nil, string, or array of strings
	if temp.Messages == nil {
		e.Messages = []string{}
	} else if msgStr, ok := temp.Messages.(string); ok {
		// Single string - convert to array
		e.Messages = []string{msgStr}
	} else if msgArray, ok := temp.Messages.([]interface{}); ok {
		// Array of strings - validate all elements are strings
		e.Messages = make([]string, 0, len(msgArray))
		for _, v := range msgArray {
			if str, ok := v.(string); ok {
				e.Messages = append(e.Messages, str)
			} else {
				return fmt.Errorf("messages array contains non-string element: %T", v)
			}
		}
	} else {
		// Unknown type - return error
		return fmt.Errorf("messages field must be nil, string, or array of strings, got %T", temp.Messages)
	}

	return nil
}

func (err APIError) String() string {
	errorString := ""

	if err.Error != "" {
		errorString += fmt.Sprintf("error: %s\n", err.Error)
	}

	if err.Detail != "" {
		errorString += fmt.Sprintf("detail: %s\n", err.Detail)
	}

	if len(err.Messages) > 0 {
		errorString += fmt.Sprintf("messages: %s", strings.Join(err.Messages, "\n"))
	}

	return errorString
}
