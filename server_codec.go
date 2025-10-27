package codec

import (
	"reflect"

	"github.com/cmd-stream/core-go"
)

// NewServerCodec creates a Codec for the server side.
//
// The server codec encodes Results and decodes Commands.
// The cmds slice lists Command types the server can handle.
// The results slice lists Result types that can be returned to the client.
func NewServerCodec[T any](cmds []reflect.Type, results []reflect.Type) (
	codec Codec[core.Result, core.Cmd[T]],
) {
	return NewCodec[core.Result, core.Cmd[T]](results, cmds)
}
