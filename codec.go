// Package codec provides a JSON-based codec implementation for cmd-stream.
package codec

import (
	"fmt"
	"reflect"

	"github.com/cmd-stream/cmd-stream-go/transport"
	cdc "github.com/cmd-stream/codec-go"
)

const ErrorPrefix = "codecjson: "

// codec represents a generic type-safe JSON codec.
// T is the type used for encoding, V is the type used for decoding.
type codec[T, V any] struct {
	cdc.Codec[T, V]
}

// newCodec constructs a JSON codec.
//
// Parameters:
//   - types1 lists the Go types that will be encoded by the Serializer.
//   - types2 lists the Go types that will be decoded by the Serializer.
func newCodec[T, V any](types1 []reflect.Type,
	types2 []reflect.Type,
) codec[T, V] {
	return codec[T, V]{
		cdc.NewCodec(types1, types2, Serializer[T, V]{}),
	}
}

// Encode writes a value of type T to the given transport.Writer.
// Returns the total number of bytes written and any error.
func (c codec[T, V]) Encode(t T, w transport.Writer) (n int, err error) {
	n, err = c.Codec.Encode(t, w)
	if err != nil {
		err = fmt.Errorf(ErrorPrefix+"%w", err)
	}
	return
}

// Decode reads a value of type V from the given transport.Reader.
// Returns the decoded value, total bytes read, and any error.
func (c codec[T, V]) Decode(r transport.Reader) (v V, n int, err error) {
	v, n, err = c.Codec.Decode(r)
	if err != nil {
		err = fmt.Errorf(ErrorPrefix+"%w", err)
	}
	return
}
