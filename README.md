# codec-json-go

[![Go Reference](https://pkg.go.dev/badge/github.com/cmd-stream/codec-json-go.svg)](https://pkg.go.dev/github.com/cmd-stream/codec-json-go)
[![GoReportCard](https://goreportcard.com/badge/cmd-stream/codec-json-go)](https://goreportcard.com/report/github.com/cmd-stream/codec-json-go)
[![codecov](https://codecov.io/gh/cmd-stream/codec-json-go/graph/badge.svg?token=nu4ycOC9bT)](https://codecov.io/gh/cmd-stream/codec-json-go)

**codec-json** provides a JSON-based codec for [cmd-stream](https://github.com/cmd-stream/cmd-stream-go).

It maps concrete Command and Result types to internal identifiers, allowing 
type-safe serialization across network boundaries.

## How To

```go
import (
  "reflect"
  cdc "github.com/cmd-stream/codec-json-go"
)

var (
  // Note: The order of types matters — two codecs created with the same types
  // in a different order are not considered equal.
  cmdTypes = []reflect.Type{
    reflect.TypeFor[YourCmd](),
    // ...
  }
  resultTypes = []reflect.Type{
    reflect.TypeFor[YourResult](),
    // ...
  }
  serverCodec = cdc.NewServerCodec(cmdTypes, resultTypes)
  clientCodec = cdc.NewClientCodec(cmdTypes, resultTypes)
)
```
