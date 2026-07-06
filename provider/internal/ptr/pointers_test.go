package ptr_test

import (
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/provider/internal/ptr"
)

func TestOf_Int(t *testing.T) {
	// Assemble
	input := 42

	// Act
	result := ptr.Of(input)

	// Assert
	if result == nil {
		t.Errorf("expected pointer to int, got nil")
	} else if *result != input {
		t.Errorf("expected %d, got %d", input, *result)
	}
}
