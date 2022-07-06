package firehydrant

import (
	"fmt"
	"strings"
)

type APIError struct {
	Error    string   `json:"error"`
	Detail   string   `json:"detail"`
	Messages []string `json:"messages"`
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
