package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		switch {
		case errors.Is(err, errHelp):
			printUsage(os.Stdout)
		case errors.Is(err, errUsage):
			fmt.Fprintln(os.Stderr, err)
			printUsage(os.Stderr)
			os.Exit(2)
		default:
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no arguments provided: %w", errUsage)
	}

	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	switch {
	case opts.view:
		if len(opts.inputs) != 1 {
			return fmt.Errorf("--view expects exactly one input file: %w", errUsage)
		}
		return viewFile(opts.inputs[0])
	case opts.batchTarget != FormatUnknown:
		if len(opts.inputs) == 0 {
			return fmt.Errorf("no input files provided for batch conversion: %w", errUsage)
		}
		return batchConvert(opts)
	case opts.stdoutFormat != FormatUnknown:
		if len(opts.inputs) != 1 {
			return fmt.Errorf("exactly one input file required when converting to stdout: %w", errUsage)
		}
		fromFormat, err := opts.resolveFromFormat(opts.inputs[0])
		if err != nil {
			return err
		}
		return convertToStdout(opts.inputs[0], fromFormat, opts.stdoutFormat)
	default:
		if len(opts.inputs) != 2 {
			return fmt.Errorf("expected input and output files: %w", errUsage)
		}
		return convertFilePair(opts.inputs[0], opts.inputs[1], opts)
	}
}

func parseArgs(args []string) (options, error) {
	var opts options

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			opts.inputs = append(opts.inputs, args[i+1:]...)
			break
		}

		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-h", "--help":
				return opts, errHelp
			case "-v", "--view":
				opts.view = true
			case "--json":
				if opts.stdoutFormat != FormatUnknown {
					return opts, fmt.Errorf("multiple stdout formats specified: %w", errUsage)
				}
				opts.stdoutFormat = FormatJSON
			case "--yaml":
				if opts.stdoutFormat != FormatUnknown {
					return opts, fmt.Errorf("multiple stdout formats specified: %w", errUsage)
				}
				opts.stdoutFormat = FormatYAML
			case "--from":
				if i+1 >= len(args) {
					return opts, fmt.Errorf("--from requires a format: %w", errUsage)
				}
				format, err := parseFormat(args[i+1])
				if err != nil {
					return opts, err
				}
				opts.from = format
				opts.hasFrom = true
				i++
			case "--to":
				if i+1 >= len(args) {
					return opts, fmt.Errorf("--to requires a format: %w", errUsage)
				}
				format, err := parseFormat(args[i+1])
				if err != nil {
					return opts, err
				}
				opts.to = format
				opts.hasTo = true
				i++
			case "--to-json":
				if opts.batchTarget != FormatUnknown {
					return opts, fmt.Errorf("multiple batch targets specified: %w", errUsage)
				}
				opts.batchTarget = FormatJSON
			case "--to-yaml":
				if opts.batchTarget != FormatUnknown {
					return opts, fmt.Errorf("multiple batch targets specified: %w", errUsage)
				}
				opts.batchTarget = FormatYAML
			case "--to-msgpack":
				if opts.batchTarget != FormatUnknown {
					return opts, fmt.Errorf("multiple batch targets specified: %w", errUsage)
				}
				opts.batchTarget = FormatMsgpack
			default:
				return opts, fmt.Errorf("unknown flag %q: %w", arg, errUsage)
			}
		} else {
			opts.inputs = append(opts.inputs, arg)
		}
	}

	if opts.view && opts.stdoutFormat != FormatUnknown {
		return opts, fmt.Errorf("--view cannot be combined with --json/--yaml: %w", errUsage)
	}
	if opts.view && opts.batchTarget != FormatUnknown {
		return opts, fmt.Errorf("--view cannot be combined with batch conversion flags: %w", errUsage)
	}
	if opts.stdoutFormat != FormatUnknown && opts.batchTarget != FormatUnknown {
		return opts, fmt.Errorf("stdout conversion cannot be combined with batch conversion flags: %w", errUsage)
	}
	if opts.hasTo && opts.batchTarget != FormatUnknown {
		return opts, fmt.Errorf("--to cannot be combined with --to-json/--to-yaml/--to-msgpack: %w", errUsage)
	}

	return opts, nil
}

func convertFilePair(inputPath, outputPath string, opts options) error {
	fromFormat, err := opts.resolveFromFormat(inputPath)
	if err != nil {
		return err
	}

	toFormat, err := opts.resolveToFormat(outputPath)
	if err != nil {
		return err
	}

	return convertAndWrite(inputPath, outputPath, fromFormat, toFormat)
}

func batchConvert(opts options) error {
	target := opts.batchTarget
	for _, input := range opts.inputs {
		fromFormat, err := opts.resolveFromFormat(input)
		if err != nil {
			return err
		}
		outputPath := deriveBatchDestination(input, target)
		if err := convertAndWrite(input, outputPath, fromFormat, target); err != nil {
			return err
		}
	}
	return nil
}

func convertAndWrite(inputPath, outputPath string, fromFormat, toFormat Format) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", inputPath, err)
	}

	converted, err := convertData(data, fromFormat, toFormat)
	if err != nil {
		return fmt.Errorf("convert %s to %s: %w", fromFormat, toFormat, err)
	}

	if needsTrailingNewline(toFormat) {
		converted = appendNewline(converted)
	}

	if err := os.WriteFile(outputPath, converted, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outputPath, err)
	}

	return nil
}

func convertToStdout(inputPath string, fromFormat, toFormat Format) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", inputPath, err)
	}

	converted, err := convertData(data, fromFormat, toFormat)
	if err != nil {
		return fmt.Errorf("convert %s to %s: %w", fromFormat, toFormat, err)
	}

	if needsTrailingNewline(toFormat) {
		converted = appendNewline(converted)
	}

	_, err = os.Stdout.Write(converted)
	return err
}

func viewFile(inputPath string) error {
	return convertToStdout(inputPath, FormatMsgpack, FormatJSON)
}

func (f Format) String() string {
	return string(f)
}

func (f Format) DefaultExt() string {
	switch f {
	case FormatMsgpack:
		return "msgpack"
	case FormatJSON:
		return "json"
	case FormatYAML:
		return "yaml"
	default:
		return ""
	}
}

func parseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "msgpack", "mpk":
		return FormatMsgpack, nil
	case "json":
		return FormatJSON, nil
	case "yaml", "yml":
		return FormatYAML, nil
	default:
		return FormatUnknown, fmt.Errorf("unknown format %q: %w", s, errUsage)
	}
}

func detectFormat(path string) (Format, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".msgpack", ".mpk":
		return FormatMsgpack, nil
	case ".json":
		return FormatJSON, nil
	case ".yaml", ".yml":
		return FormatYAML, nil
	default:
		return FormatUnknown, fmt.Errorf("unable to infer format from %q: %w", path, errUsage)
	}
}

func deriveBatchDestination(inputPath string, target Format) string {
	base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
	ext := target.DefaultExt()
	if ext == "" {
		return inputPath
	}
	return base + "." + ext
}

func (o options) resolveFromFormat(path string) (Format, error) {
	if o.hasFrom {
		return o.from, nil
	}
	return detectFormat(path)
}

func (o options) resolveToFormat(path string) (Format, error) {
	if o.hasTo {
		return o.to, nil
	}
	return detectFormat(path)
}

func convertData(data []byte, fromFormat, toFormat Format) ([]byte, error) {
	if fromFormat == FormatUnknown || toFormat == FormatUnknown {
		return nil, fmt.Errorf("unsupported conversion from %q to %q", fromFormat, toFormat)
	}

	value, err := decodeData(data, fromFormat)
	if err != nil {
		return nil, err
	}

	return encodeData(value, toFormat)
}

func decodeData(data []byte, format Format) (interface{}, error) {
	var value interface{}
	switch format {
	case FormatMsgpack:
		if err := msgpack.Unmarshal(data, &value); err != nil {
			return nil, fmt.Errorf("decode msgpack: %w", err)
		}
	case FormatJSON:
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&value); err != nil {
			return nil, fmt.Errorf("decode json: %w", err)
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &value); err != nil {
			return nil, fmt.Errorf("decode yaml: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
	return normalizeValue(value), nil
}

func encodeData(value interface{}, format Format) ([]byte, error) {
	switch format {
	case FormatMsgpack:
		return msgpack.Marshal(value)
	case FormatJSON:
		return json.MarshalIndent(value, "", "  ")
	case FormatYAML:
		return yaml.Marshal(value)
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}
}

func normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		out := make(map[string]interface{}, len(v))
		for key, val := range v {
			out[fmt.Sprint(key)] = normalizeValue(val)
		}
		return out
	case map[string]interface{}:
		for key, val := range v {
			v[key] = normalizeValue(val)
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = normalizeValue(val)
		}
		return v
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
		}
		if f, err := v.Float64(); err == nil {
			return f
		}
		return v.String()
	default:
		return v
	}
}

func needsTrailingNewline(format Format) bool {
	return format == FormatJSON || format == FormatYAML
}

func appendNewline(data []byte) []byte {
	if len(data) == 0 || data[len(data)-1] != '\n' {
		return append(data, '\n')
	}
	return data
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "mpt v"+versionText)
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  mpt --view file.msgpack")
	fmt.Fprintln(w, "  mpt input.msgpack output.json")
	fmt.Fprintln(w, "  mpt --from msgpack --to json input.bin output.txt")
	fmt.Fprintln(w, "  mpt data.msgpack --json")
	fmt.Fprintln(w, "  mpt *.msgpack --to-json")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "options:")
	fmt.Fprintln(w, "  -h, --help          show this help message")
	fmt.Fprintln(w, "  -v, --view          render messagepack as json to stdout")
	fmt.Fprintln(w, "      --json          convert input to json and write to stdout")
	fmt.Fprintln(w, "      --yaml          convert input to yaml and write to stdout")
	fmt.Fprintln(w, "      --from format   override detected input format")
	fmt.Fprintln(w, "      --to format     override detected output format for single conversion")
	fmt.Fprintln(w, "      --to-json       batch convert input files to json files")
	fmt.Fprintln(w, "      --to-yaml       batch convert input files to yaml files")
	fmt.Fprintln(w, "      --to-msgpack    batch convert input files to messagepack files")
}
