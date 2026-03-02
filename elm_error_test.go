package elmobd

import (
	"errors"
	"strings"
	"testing"
)

// TestNewResultELMErrors verifies that NewResult detects known ELM327 textual
// responses and returns an *ELMError instead of a generic parse error.
func TestNewResultELMErrors(t *testing.T) {
	elmLines := []string{
		"CAN ERROR",
		"NO DATA",
		"STOPPED",
		"UNABLE TO CONNECT",
		"SEARCHING...",
		"BUS INIT",
		// Variations with leading/trailing whitespace
		"  CAN ERROR  ",
		"  NO DATA",
		// Mixed case
		"can error",
		"no data",
	}

	for _, line := range elmLines {
		_, err := NewResult(line)

		if err == nil {
			t.Errorf("Expected error for ELM line %q, got nil", line)
			continue
		}

		var elmErr *ELMError
		if !errors.As(err, &elmErr) {
			t.Errorf("Expected *ELMError for ELM line %q, got %T: %v", line, err, err)
		}
	}
}

// TestNewResultValidHex verifies that a well-formed hex response still parses
// correctly after the ELM detection logic was added.
func TestNewResultValidHex(t *testing.T) {
	result, err := NewResult("41 0C FF B2")

	assertSuccess(t, err)

	if result == nil {
		t.Fatal("Expected non-nil result for valid hex response")
	}
}

// TestNewResultMultipleSpaces verifies that multiple spaces between hex bytes
// are handled correctly (using strings.Fields instead of strings.Split).
func TestNewResultMultipleSpaces(t *testing.T) {
	result, err := NewResult("41  0C  FF  B2")

	assertSuccess(t, err)

	if result == nil {
		t.Fatal("Expected non-nil result for hex response with multiple spaces")
	}
}

// TestParseOBDResponseELMErrors verifies that parseOBDResponse returns
// *ELMError for known ELM327 status/error lines.
func TestParseOBDResponseELMErrors(t *testing.T) {
	type scenario struct {
		outputs []string
	}

	scenarios := []scenario{
		{[]string{"CAN ERROR"}},
		{[]string{"NO DATA"}},
		{[]string{"STOPPED"}},
		{[]string{"UNABLE TO CONNECT"}},
		// ELM error after a SEARCHING line
		{[]string{"SEARCHING...", "CAN ERROR"}},
	}

	cmd := NewEngineLoad()

	for _, curr := range scenarios {
		_, err := parseOBDResponse(cmd, curr.outputs)

		if err == nil {
			t.Errorf("Expected error for outputs %v, got nil", curr.outputs)
			continue
		}

		var elmErr *ELMError
		if !errors.As(err, &elmErr) {
			t.Errorf("Expected *ELMError for outputs %v, got %T: %v", curr.outputs, err, err)
		}
	}
}

// TestELMErrorMessage verifies that ELMError.Error() includes the original raw line.
func TestELMErrorMessage(t *testing.T) {
	raw := "CAN ERROR"
	err := &ELMError{Message: raw}

	msg := err.Error()
	if msg == "" {
		t.Fatal("Expected non-empty error message")
	}

	// The original raw line should appear in the error message.
	if !strings.Contains(msg, raw) {
		t.Errorf("Expected error message to contain %q, got %q", raw, msg)
	}
}
