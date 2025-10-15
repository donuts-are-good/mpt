package main

import (
	"path/filepath"
	"testing"
)

func TestYAMLComments(t *testing.T) {
	dir := setupTestDir(t)

	yamlWithComments := []byte(`# This is a comment
key: value
# Another comment
nested:
  # Nested comment
  data: test`)

	input := filepath.Join(dir, "comments.yaml")
	output := filepath.Join(dir, "output.json")

	writeTestFile(t, input, yamlWithComments)

	testConvertFile(t, input, output, FormatYAML, FormatJSON)

	jsonData, _ := readFile(output)
	assertValidJSON(t, jsonData)
}

func TestJSONPrecision(t *testing.T) {
	dir := setupTestDir(t)

	preciseJSON := []byte(`{
		"maxInt": 9007199254740991,
		"minInt": -9007199254740991,
		"float": 3.141592653589793,
		"verySmall": 0.0000000000001,
		"veryLarge": 999999999999.999
	}`)

	input := filepath.Join(dir, "precision.json")
	msgpack := filepath.Join(dir, "precision.msgpack")
	roundtrip := filepath.Join(dir, "precision-rt.json")

	writeTestFile(t, input, preciseJSON)

	testConvertFile(t, input, msgpack, FormatJSON, FormatMsgpack)
	testConvertFile(t, msgpack, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertValidJSON(t, rtData)
}

func TestMsgpackBinaryVsString(t *testing.T) {
	dir := setupTestDir(t)

	testJSON := []byte(`{"text":"hello","bytes":"base64data","number":42}`)

	input := filepath.Join(dir, "binary.json")
	msgpack := filepath.Join(dir, "binary.msgpack")
	roundtrip := filepath.Join(dir, "binary-rt.json")

	writeTestFile(t, input, testJSON)

	testConvertFile(t, input, msgpack, FormatJSON, FormatMsgpack)

	mpData, _ := readFile(msgpack)
	assertValidMsgpack(t, mpData)

	testConvertFile(t, msgpack, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertJSONEqual(t, testJSON, rtData)
}
