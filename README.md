# codec-json

**codec-json** provides a JSON-based codec for [cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go).

It maps concrete Command and Result types to internal identifiers,
allowing type-safe serialization across network boundaries.

## Example

```go
var (
  cmdTypes = []reflect.Type{
    reflect.TypeFor[YourCmd](),
    // ...
  }
  resultTypes = []reflect.Type{
    reflect.TypeFor[YourResult](),
    // ...
  }

  serverCodec = codecjson.NewServerCodec(cmdTypes, resultTypes)
  clientCodec = codecjson.NewClientCodec(cmdTypes, resultTypes)
)
```
