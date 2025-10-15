package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vmihailenco/msgpack/v5"

	"gopkg.in/yaml.v3"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "mpt-test-*")
	if err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

func writeTestFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write test file %s: %v", path, err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", path)
	}
}

func assertValidJSON(t *testing.T, data []byte) {
	t.Helper()
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		t.Errorf("invalid json: %v", err)
	}
}

func assertValidYAML(t *testing.T, data []byte) {
	t.Helper()
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		t.Errorf("invalid yaml: %v", err)
	}
}

func assertValidMsgpack(t *testing.T, data []byte) {
	t.Helper()
	var v interface{}
	if err := msgpack.Unmarshal(data, &v); err != nil {
		t.Errorf("invalid msgpack: %v", err)
	}
}

func assertJSONEqual(t *testing.T, expected, actual []byte) {
	t.Helper()

	var exp, act interface{}
	if err := json.Unmarshal(expected, &exp); err != nil {
		t.Fatalf("expected is not valid json: %v", err)
	}
	if err := json.Unmarshal(actual, &act); err != nil {
		t.Fatalf("actual is not valid json: %v", err)
	}

	expNorm, _ := json.Marshal(exp)
	actNorm, _ := json.Marshal(act)

	if !bytes.Equal(expNorm, actNorm) {
		t.Errorf("json mismatch:\nexpected: %s\nactual: %s", expNorm, actNorm)
	}
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("content", name))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return data
}

func testConvertFile(t *testing.T, inputPath, outputPath string, fromFormat, toFormat Format) {
	t.Helper()
	if err := convertAndWrite(inputPath, outputPath, fromFormat, toFormat); err != nil {
		t.Fatalf("convertAndWrite failed: %v", err)
	}
}

func assertError(t *testing.T, err error, msgContains string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", msgContains)
		return
	}
	if msgContains != "" && !strings.Contains(err.Error(), msgContains) {
		t.Errorf("expected error containing %q, got: %v", msgContains, err)
	}
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
