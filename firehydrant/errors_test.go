package firehydrant

import (
	"encoding/json"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("all attributes", func(t *testing.T) {
		jsonData := []byte(`{"error":"test error","detail":"test detail","messages":["test message1","test message2"]}`)
		apiError := &APIError{}

		err := json.Unmarshal(jsonData, apiError)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		expected := &APIError{
			Error:    "test error",
			Detail:   "test detail",
			Messages: []string{"test message1", "test message2"},
		}
		got := apiError
		if expected.Error != got.Error {
			t.Errorf("unexpected Error\nexpected: %s\ngot: %s", expected.Error, got.Error)
		}
		if expected.Detail != got.Detail {
			t.Errorf("unexpected Detail\nexpected: %s\ngot: %s", expected.Detail, got.Detail)
		}

		if len(expected.Messages) != len(got.Messages) {
			t.Errorf("unexpected Messages lengtth\nexpected: %d\ngot: %d", len(expected.Messages), len(got.Messages))
		}

		for index := range expected.Messages {
			if expected.Messages[index] != got.Messages[index] {
				t.Errorf("unexpected Messages[%d]\nexpected: %s\ngot: %s", index, expected.Messages[index], got.Messages[index])
			}
		}
	})

	t.Run("no attributes", func(t *testing.T) {
		jsonData := []byte(`{}`)
		apiError := &APIError{}

		err := json.Unmarshal(jsonData, apiError)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		expected := &APIError{}
		got := apiError
		if expected.Error != got.Error {
			t.Errorf("unexpected Error\nexpected: %s\ngot: %s", expected.Error, got.Error)
		}
		if expected.Detail != got.Detail {
			t.Errorf("unexpected Detail\nexpected: %s\ngot: %s", expected.Detail, got.Detail)
		}

		if len(expected.Messages) != len(got.Messages) {
			t.Errorf("unexpected Messages length\nexpected: %d\ngot: %d", len(expected.Messages), len(got.Messages))
		}
	})

	t.Run("extra attributes", func(t *testing.T) {
		jsonData := []byte(`{"extra": 1,"error":"test error","detail":"test detail","messages":["test message1","test message2"]}`)
		apiError := &APIError{}

		err := json.Unmarshal(jsonData, apiError)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		expected := &APIError{
			Error:    "test error",
			Detail:   "test detail",
			Messages: []string{"test message1", "test message2"},
		}
		got := apiError
		if expected.Error != got.Error {
			t.Errorf("unexpected Error\nexpected: %s\ngot: %s", expected.Error, got.Error)
		}
		if expected.Detail != got.Detail {
			t.Errorf("unexpected Detail\nexpected: %s\ngot: %s", expected.Detail, got.Detail)
		}

		if len(expected.Messages) != len(got.Messages) {
			t.Errorf("unexpected Messages length\nexpected: %d\ngot: %d", len(expected.Messages), len(got.Messages))
		}
		if len(expected.Messages) != len(got.Messages) {
			t.Errorf("unexpected Messages length\nexpected: %d\ngot: %d", len(expected.Messages), len(got.Messages))
		}

		for index := range expected.Messages {
			if expected.Messages[index] != got.Messages[index] {
				t.Errorf("unexpected Messages[%d]\nexpected: %s\ngot: %s", index, expected.Messages[index], got.Messages[index])
			}
		}
	})
}

func TestErrors_String(t *testing.T) {
	t.Run("all attributes", func(t *testing.T) {
		apiError := &APIError{
			Error:    "test error",
			Detail:   "test detail",
			Messages: []string{"test message1", "test message2"},
		}

		expected := "error: test error\ndetail: test detail\nmessages: test message1\ntest message2"
		got := apiError.String()
		if expected != got {
			t.Errorf("unexpected string\nexpected: %s\ngot: %s", expected, got)
		}
	})

	t.Run("no attributes", func(t *testing.T) {
		apiError := &APIError{}

		expected := ""
		got := apiError.String()
		if expected != got {
			t.Errorf("unexpected string\nexpected empty string\ngot: %s", got)
		}
	})

	t.Run("no detail", func(t *testing.T) {
		apiError := &APIError{
			Error:    "test error",
			Messages: []string{"test message1", "test message2"},
		}

		expected := "error: test error\nmessages: test message1\ntest message2"
		got := apiError.String()
		if expected != got {
			t.Errorf("unexpected string\nexpected: %s\ngot: %s", expected, got)
		}
	})

	t.Run("no error", func(t *testing.T) {
		apiError := &APIError{
			Detail:   "test detail",
			Messages: []string{"test message1", "test message2"},
		}

		expected := "detail: test detail\nmessages: test message1\ntest message2"
		got := apiError.String()
		if expected != got {
			t.Errorf("unexpected string\nexpected: %s\ngot: %s", expected, got)
		}
	})

	t.Run("no messages", func(t *testing.T) {
		apiError := &APIError{
			Error:    "test error",
			Detail:   "test detail",
			Messages: []string{},
		}

		expected := "error: test error\ndetail: test detail\n"
		got := apiError.String()
		if expected != got {
			t.Errorf("unexpected string\nexpected: %s\ngot: %s", expected, got)
		}
	})
}
