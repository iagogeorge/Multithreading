package main

import (
	"bytes"
	"log"
	"testing"
)

// TestGetAddress function to test the getAddress function
func TestGetAddress(t *testing.T) {
	// Capturing log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)

	cep := "01153000"
	getAddress(cep)

	// Resetting log output
	log.SetOutput(nil)

	logOutput := logBuf.String()

	if !contains(logOutput, "Fastest API:") {
		t.Error("Expected log output to contain 'Fastest API:'")
	}

	if !contains(logOutput, "Address:") {
		t.Error("Expected log output to contain 'Address:'")
	}

	if contains(logOutput, "Error") {
		t.Error("Unexpected error in log output")
	}
}

// Helper function to check if the log output contains a specific substring
func contains(logOutput, substring string) bool {
	return bytes.Contains([]byte(logOutput), []byte(substring))
}
