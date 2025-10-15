package main

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestNoArguments(t *testing.T) {
	err := run([]string{})
	assertError(t, err, "no arguments provided")
	if !errors.Is(err, errUsage) {
		t.Errorf("expected errUsage, got: %v", err)
	}
}

func TestInvalidFlags(t *testing.T) {
	err := run([]string{"--invalid-flag"})
	assertError(t, err, "unknown flag")
	if !errors.Is(err, errUsage) {
		t.Errorf("expected errUsage, got: %v", err)
	}
}

func TestHelpUsage(t *testing.T) {
	err := run([]string{"-h"})
	if !errors.Is(err, errHelp) {
		t.Errorf("expected errHelp, got: %v", err)
	}

	err = run([]string{"--help"})
	if !errors.Is(err, errHelp) {
		t.Errorf("expected errHelp, got: %v", err)
	}
}

func TestConflictingFlags(t *testing.T) {
	_, err := parseArgs([]string{"--view", "--json", "file.json"})
	assertError(t, err, "cannot be combined")

	_, err = parseArgs([]string{"--to-json", "--to-yaml", "file.msgpack"})
	assertError(t, err, "multiple batch targets")
}

func TestFilePathEdgeCases(t *testing.T) {
	dir := setupTestDir(t)

	testData := []byte(`{"test":"value"}`)
	spaceFile := filepath.Join(dir, "file with spaces.json")
	unicodeFile := filepath.Join(dir, "файл-文件.json")
	outputSpace := filepath.Join(dir, "output with spaces.msgpack")
	outputUnicode := filepath.Join(dir, "输出-файл.msgpack")

	writeTestFile(t, spaceFile, testData)
	writeTestFile(t, unicodeFile, testData)

	testConvertFile(t, spaceFile, outputSpace, FormatJSON, FormatMsgpack)
	assertFileExists(t, outputSpace)

	testConvertFile(t, unicodeFile, outputUnicode, FormatJSON, FormatMsgpack)
	assertFileExists(t, outputUnicode)
}
