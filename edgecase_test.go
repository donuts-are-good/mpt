package main

import (
	"path/filepath"
	"testing"
)

func TestEmptyFiles(t *testing.T) {
	dir := setupTestDir(t)

	emptyFile := filepath.Join(dir, "empty.json")
	output := filepath.Join(dir, "output.msgpack")

	writeTestFile(t, emptyFile, []byte(""))

	err := convertAndWrite(emptyFile, output, FormatJSON, FormatMsgpack)
	assertError(t, err, "decode")
}

func TestEmptyObjects(t *testing.T) {
	dir := setupTestDir(t)

	emptyObj := []byte("{}")
	input := filepath.Join(dir, "empty.json")
	output := filepath.Join(dir, "output.msgpack")
	roundtrip := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, input, emptyObj)

	testConvertFile(t, input, output, FormatJSON, FormatMsgpack)
	testConvertFile(t, output, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertJSONEqual(t, emptyObj, rtData)
}

func TestEmptyArrays(t *testing.T) {
	dir := setupTestDir(t)

	emptyArr := []byte("[]")
	input := filepath.Join(dir, "empty.json")
	output := filepath.Join(dir, "output.msgpack")
	roundtrip := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, input, emptyArr)

	testConvertFile(t, input, output, FormatJSON, FormatMsgpack)
	testConvertFile(t, output, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertJSONEqual(t, emptyArr, rtData)
}

func TestNullValues(t *testing.T) {
	dir := setupTestDir(t)

	nullData := []byte(`{"value":null,"nested":{"also":null}}`)
	input := filepath.Join(dir, "null.json")
	output := filepath.Join(dir, "output.msgpack")
	roundtrip := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, input, nullData)

	testConvertFile(t, input, output, FormatJSON, FormatMsgpack)
	testConvertFile(t, output, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertJSONEqual(t, nullData, rtData)
}

func TestLargeNumbers(t *testing.T) {
	dir := setupTestDir(t)

	numData := []byte(`{"int":9223372036854775807,"float":1.7976931348623157e+308,"scientific":1.23e-10}`)
	input := filepath.Join(dir, "numbers.json")
	output := filepath.Join(dir, "output.msgpack")
	roundtrip := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, input, numData)

	testConvertFile(t, input, output, FormatJSON, FormatMsgpack)
	testConvertFile(t, output, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertValidJSON(t, rtData)
}

func TestUnicodeChars(t *testing.T) {
	dir := setupTestDir(t)

	unicodeData := []byte(`{"emoji":"🔥💯","chinese":"你好世界","arabic":"مرحبا","special":"¡Hola!"}`)
	input := filepath.Join(dir, "unicode.json")
	output := filepath.Join(dir, "output.msgpack")
	roundtrip := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, input, unicodeData)

	testConvertFile(t, input, output, FormatJSON, FormatMsgpack)
	testConvertFile(t, output, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertJSONEqual(t, unicodeData, rtData)
}

func TestDeeplyNested(t *testing.T) {
	dir := setupTestDir(t)

	deepData := []byte(`{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":"deep"}}}}}}}}}}`)
	input := filepath.Join(dir, "deep.json")
	output := filepath.Join(dir, "output.msgpack")
	roundtrip := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, input, deepData)

	testConvertFile(t, input, output, FormatJSON, FormatMsgpack)
	testConvertFile(t, output, roundtrip, FormatMsgpack, FormatJSON)

	rtData, _ := readFile(roundtrip)
	assertJSONEqual(t, deepData, rtData)
}
