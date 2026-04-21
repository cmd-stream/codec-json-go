package codec

import (
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// ServerCodec defines a JSON codec for the server side.
type ServerCodec[T any] = codec[core.Result, core.Cmd[T]]

// NewServerCodec creates a JSON codec for the server side.
//
// The cmdTypes slice lists Command types the server can handle.
// The resultTypes slice lists Result types that can be returned to the client.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
func NewServerCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	c ServerCodec[T],
) {
	return newCodec[core.Result, core.Cmd[T]](resultTypes, cmdTypes)
}

// NewServerCodecWith creates a JSON codec for the server side using the
// provided Registry.
func NewServerCodecWith[T any](registry *Registry[T]) ServerCodec[T] {
	return NewServerCodec[T](registry.Cmds(), registry.Results())
}

