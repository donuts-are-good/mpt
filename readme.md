# mpt
message pack tools

## syntax

view msgpack as plaintext
```
mpt --view something.msgpack
mpt -v something.msgpack
```

```
example data from that mpt view example
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

## license
2025 mit license