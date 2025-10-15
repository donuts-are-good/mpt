package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJSONMsgpackRoundtrip(t *testing.T) {
	dir := setupTestDir(t)

	jsonInput := loadFixture(t, "json/demo1.json")

	jsonPath := filepath.Join(dir, "input.json")
	msgpackPath := filepath.Join(dir, "output.msgpack")
	roundtripPath := filepath.Join(dir, "roundtrip.json")

	writeTestFile(t, jsonPath, jsonInput)

	testConvertFile(t, jsonPath, msgpackPath, FormatJSON, FormatMsgpack)
	assertFileExists(t, msgpackPath)

	msgpackData, _ := os.ReadFile(msgpackPath)
	assertValidMsgpack(t, msgpackData)

	testConvertFile(t, msgpackPath, roundtripPath, FormatMsgpack, FormatJSON)
	assertFileExists(t, roundtripPath)

	roundtripData, err := os.ReadFile(roundtripPath)
	if err != nil {
		t.Fatalf("failed to read roundtrip file: %v", err)
	}

	assertJSONEqual(t, jsonInput, roundtripData)
}

func TestViewMsgpack(t *testing.T) {
	dir := setupTestDir(t)

	jsonInput := loadFixture(t, "json/demo1.json")

	jsonPath := filepath.Join(dir, "input.json")
	msgpackPath := filepath.Join(dir, "input.msgpack")

	writeTestFile(t, jsonPath, jsonInput)
	testConvertFile(t, jsonPath, msgpackPath, FormatJSON, FormatMsgpack)

	viewOutput, err := readAndConvert(msgpackPath, FormatMsgpack, FormatJSON)
	if err != nil {
		t.Fatalf("viewFile failed: %v", err)
	}

	assertValidJSON(t, viewOutput)
	assertJSONEqual(t, jsonInput, viewOutput)
}

func TestViewJSON(t *testing.T) {
	dir := setupTestDir(t)

	jsonInput := loadFixture(t, "json/demo1.json")
	jsonPath := filepath.Join(dir, "input.json")

	writeTestFile(t, jsonPath, jsonInput)

	viewOutput, err := readAndConvert(jsonPath, FormatJSON, FormatJSON)
	if err != nil {
		t.Fatalf("viewFile on json failed: %v", err)
	}

	assertValidJSON(t, viewOutput)
	assertJSONEqual(t, jsonInput, viewOutput)
}

func TestYAMLRoundtrip(t *testing.T) {
	dir := setupTestDir(t)

	yamlInput := loadFixture(t, "yaml/demo1.yaml")

	yamlPath := filepath.Join(dir, "input.yaml")
	msgpackPath := filepath.Join(dir, "output.msgpack")
	roundtripPath := filepath.Join(dir, "roundtrip.yaml")

	writeTestFile(t, yamlPath, yamlInput)

	testConvertFile(t, yamlPath, msgpackPath, FormatYAML, FormatMsgpack)
	assertFileExists(t, msgpackPath)

	msgpackData, _ := os.ReadFile(msgpackPath)
	assertValidMsgpack(t, msgpackData)

	testConvertFile(t, msgpackPath, roundtripPath, FormatMsgpack, FormatYAML)
	assertFileExists(t, roundtripPath)

	roundtripData, err := os.ReadFile(roundtripPath)
	if err != nil {
		t.Fatalf("failed to read roundtrip file: %v", err)
	}

	assertValidYAML(t, roundtripData)
}

func TestFormatOverride(t *testing.T) {
	dir := setupTestDir(t)

	jsonInput := loadFixture(t, "json/demo2.json")

	weirdInput := filepath.Join(dir, "weird.bin")
	weirdOutput := filepath.Join(dir, "weird-output.bin")
	finalOutput := filepath.Join(dir, "final.json")

	writeTestFile(t, weirdInput, jsonInput)

	testConvertFile(t, weirdInput, weirdOutput, FormatJSON, FormatMsgpack)
	assertFileExists(t, weirdOutput)

	msgpackData, _ := os.ReadFile(weirdOutput)
	assertValidMsgpack(t, msgpackData)

	testConvertFile(t, weirdOutput, finalOutput, FormatMsgpack, FormatJSON)
	assertFileExists(t, finalOutput)

	finalData, _ := os.ReadFile(finalOutput)
	assertJSONEqual(t, jsonInput, finalData)
}

func TestStdoutJSON(t *testing.T) {
	dir := setupTestDir(t)

	jsonInput := loadFixture(t, "json/demo3.json")
	jsonPath := filepath.Join(dir, "input.json")

	writeTestFile(t, jsonPath, jsonInput)

	output, err := readAndConvert(jsonPath, FormatJSON, FormatJSON)
	if err != nil {
		t.Fatalf("stdout json conversion failed: %v", err)
	}

	assertValidJSON(t, output)
	assertJSONEqual(t, jsonInput, output)
}

func TestStdoutYAML(t *testing.T) {
	dir := setupTestDir(t)

	jsonInput := loadFixture(t, "json/demo4.json")
	jsonPath := filepath.Join(dir, "input.json")

	writeTestFile(t, jsonPath, jsonInput)

	output, err := readAndConvert(jsonPath, FormatJSON, FormatYAML)
	if err != nil {
		t.Fatalf("stdout yaml conversion failed: %v", err)
	}

	assertValidYAML(t, output)
}

func TestBatchMsgpack(t *testing.T) {
	dir := setupTestDir(t)

	json0 := loadFixture(t, "json/demo0.json")
	json1 := loadFixture(t, "json/demo1.json")

	json0Path := filepath.Join(dir, "demo0.json")
	json1Path := filepath.Join(dir, "demo1.json")

	writeTestFile(t, json0Path, json0)
	writeTestFile(t, json1Path, json1)

	testConvertFile(t, json0Path, filepath.Join(dir, "demo0.msgpack"), FormatJSON, FormatMsgpack)
	testConvertFile(t, json1Path, filepath.Join(dir, "demo1.msgpack"), FormatJSON, FormatMsgpack)

	msgpack0Path := filepath.Join(dir, "demo0.msgpack")
	msgpack1Path := filepath.Join(dir, "demo1.msgpack")

	assertFileExists(t, msgpack0Path)
	assertFileExists(t, msgpack1Path)

	mp0Data, _ := os.ReadFile(msgpack0Path)
	mp1Data, _ := os.ReadFile(msgpack1Path)

	assertValidMsgpack(t, mp0Data)
	assertValidMsgpack(t, mp1Data)
}

func TestBatchJSON(t *testing.T) {
	dir := setupTestDir(t)

	json0 := loadFixture(t, "json/demo0.json")
	json1 := loadFixture(t, "json/demo1.json")

	json0Path := filepath.Join(dir, "demo0.json")
	json1Path := filepath.Join(dir, "demo1.json")

	writeTestFile(t, json0Path, json0)
	writeTestFile(t, json1Path, json1)

	msgpack0 := filepath.Join(dir, "demo0.msgpack")
	msgpack1 := filepath.Join(dir, "demo1.msgpack")

	testConvertFile(t, json0Path, msgpack0, FormatJSON, FormatMsgpack)
	testConvertFile(t, json1Path, msgpack1, FormatJSON, FormatMsgpack)

	roundtrip0 := filepath.Join(dir, "roundtrip0.json")
	roundtrip1 := filepath.Join(dir, "roundtrip1.json")

	testConvertFile(t, msgpack0, roundtrip0, FormatMsgpack, FormatJSON)
	testConvertFile(t, msgpack1, roundtrip1, FormatMsgpack, FormatJSON)

	assertFileExists(t, roundtrip0)
	assertFileExists(t, roundtrip1)

	rt0Data, _ := os.ReadFile(roundtrip0)
	rt1Data, _ := os.ReadFile(roundtrip1)

	assertJSONEqual(t, json0, rt0Data)
	assertJSONEqual(t, json1, rt1Data)
}

func TestBatchYAML(t *testing.T) {
	dir := setupTestDir(t)

	json0 := loadFixture(t, "json/demo0.json")
	json1 := loadFixture(t, "json/demo1.json")

	msgpack0 := filepath.Join(dir, "demo0.msgpack")
	msgpack1 := filepath.Join(dir, "demo1.msgpack")

	data0, _ := convertData(json0, FormatJSON, FormatMsgpack)
	data1, _ := convertData(json1, FormatJSON, FormatMsgpack)

	writeTestFile(t, msgpack0, data0)
	writeTestFile(t, msgpack1, data1)

	yaml0Path := filepath.Join(dir, "demo0.yaml")
	yaml1Path := filepath.Join(dir, "demo1.yaml")

	testConvertFile(t, msgpack0, yaml0Path, FormatMsgpack, FormatYAML)
	testConvertFile(t, msgpack1, yaml1Path, FormatMsgpack, FormatYAML)

	assertFileExists(t, yaml0Path)
	assertFileExists(t, yaml1Path)

	yaml0Data, _ := os.ReadFile(yaml0Path)
	yaml1Data, _ := os.ReadFile(yaml1Path)

	assertValidYAML(t, yaml0Data)
	assertValidYAML(t, yaml1Data)
}
