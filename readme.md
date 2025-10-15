# mpt
message pack tools

## syntax
view msgpack as plaintext
```
mpt --view something.msgpack
mpt -v something.msgpack
```

## infer implied conversions by default

smart-convert msgpack <-> json
```
mpt convert-this.msgpack into-this.json
mpt convert-this.json into-this.msgpack
```

smart-convert msgpack <-> yaml
```
mpt convert-this.yml into-this.msgpack
mpt convert-this.yaml into-this.msgpack
mpt convert-this.msgpack into-this.yml
mpt convert-this.msgpack into-this.yaml
```

## format override
use arbitrary extensions
```
mpt --from msgpack --to json input output
```

## output to stdout
output to stdout
```
mpt data.msgpack --json
mpt data.msgpack --yaml
```

## multiple file conversion
batch convert
```
mpt *.msgpack --to-json
mpt *.msgpack --to-yaml
mpt *.json --to-msgpack
mpt *.yml --to-msgpack
mpt *.yaml --to-msgpack
```

## license
2025 mit license