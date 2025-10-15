package main

import (
	"errors"
)

var (
	errUsage    = errors.New("usage error")
	errHelp     = errors.New("help requested")
	versionText = "0.0.1"
)

type options struct {
	view         bool
	stdoutFormat Format
	batchTarget  Format
	from         Format
	to           Format
	hasFrom      bool
	hasTo        bool
	inputs       []string
}

type Format string

const (
	FormatUnknown Format = ""
	FormatMsgpack Format = "msgpack"
	FormatJSON    Format = "json"
	FormatYAML    Format = "yaml"
)
