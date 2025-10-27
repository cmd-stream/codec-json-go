package codec

import (
	"reflect"

	"github.com/cmd-stream/core-go"
)

// NewClientCodec creates a Codec for the client side.
//
// The client codec encodes Commands and decodes Results.
// The cmds slice lists Command types the client can send.
// The results slice lists Result types the client expects to receive.
func NewClientCodec[T any](cmds []reflect.Type, results []reflect.Type) (
	codec Codec[core.Cmd[T], core.Result],
) {
	return NewCodec[core.Cmd[T], core.Result](cmds, results)
}
