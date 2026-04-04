package codec

import (
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/core"
)

// ClientCodec defines a JSON codec for the client side.
type ClientCodec[T any] = codec[core.Cmd[T], core.Result]

// NewClientCodec creates a JSON codec for the client side.
//
// The cmdTypes slice lists Command types the client can send.
// The resultTypes slice lists Result types the client expects to receive.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
func NewClientCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	c ClientCodec[T],
) {
	return newCodec[core.Cmd[T], core.Result](cmdTypes, resultTypes)
}
