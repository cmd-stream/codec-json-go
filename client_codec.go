package codec

import (
	"reflect"

	"github.com/cmd-stream/core-go"
)

// NewClientCodec creates a JSON Codec for the client side.
//
// The cmdTypes slice lists Command types the client can send.
// The resultTypes slice lists Result types the client expects to receive.
//
// Note: The order of types matters — two codecs created with the same types
// in a different order are not considered equal.
func NewClientCodec[T any](cmdTypes []reflect.Type, resultTypes []reflect.Type) (
	codec Codec[core.Cmd[T], core.Result],
) {
	return NewCodec[core.Cmd[T], core.Result](cmdTypes, resultTypes)
}
