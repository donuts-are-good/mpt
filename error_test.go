package main

import (
	"path/filepath"
	"testing"
)

func TestInvalidFiles(t *testing.T) {
	dir := setupTestDir(t)

	invalidJSON := filepath.Join(dir, "invalid.json")
	invalidYAML := filepath.Join(dir, "invalid.yaml")
	output := filepath.Join(dir, "output.json")

	writeTestFile(t, invalidJSON, []byte("not valid json{{{"))
	writeTestFile(t, invalidYAML, []byte("not: valid: yaml: [[["))

	err := convertAndWrite(invalidJSON, output, FormatJSON, FormatJSON)
	assertError(t, err, "decode json")

	err = convertAndWrite(invalidYAML, output, FormatYAML, FormatJSON)
	assertError(t, err, "decode yaml")
}

func TestMissingFiles(t *testing.T) {
	dir := setupTestDir(t)

	missingInput := filepath.Join(dir, "does-not-exist.json")
	output := filepath.Join(dir, "output.json")

	err := convertAndWrite(missingInput, output, FormatJSON, FormatJSON)
	assertError(t, err, "read")
}

func TestUnsupportedFormats(t *testing.T) {
	dir := setupTestDir(t)

	input := filepath.Join(dir, "file.unknown")
	output := filepath.Join(dir, "output.json")

	writeTestFile(t, input, []byte("{}"))

	_, err := detectFormat(input)
	assertError(t, err, "unable to infer format")

	err = convertAndWrite(input, output, FormatUnknown, FormatJSON)
	assertError(t, err, "unsupported conversion")
}
