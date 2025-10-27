// Package codec provides a JSON-based codec implementation for cmd-stream-go.
package codec

import (
	"encoding/json"
	"reflect"

	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
	"github.com/mus-format/dts-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
)

// NewCodec constructs a Codec with explicit type mappings.
//
// The first slice lists types that can be encoded,
// and the second slice lists types that can be decoded.
func NewCodec[T, V any](types1 []reflect.Type, types2 []reflect.Type) (
	codec Codec[T, V],
) {
	if len(types1) == 0 {
		panic(errorPrefix + "types1 is empty")
	}
	if len(types2) == 0 {
		panic(errorPrefix + "types2 is empty")
	}
	codec = Codec[T, V]{
		typeMap: make(map[reflect.Type]com.DTM),
		dtmSl:   make([]reflect.Type, len(types2)),
	}
	for i, t := range types1 {
		codec.typeMap[t] = com.DTM(i)
	}
	copy(codec.dtmSl, types2)
	return
}

// Codec provides JSON serialization and deserialization for registered types.
//
// T represents the types that can be encoded, and V represents the types that
// can be decoded.
type Codec[T, V any] struct {
	typeMap map[reflect.Type]com.DTM
	dtmSl   []reflect.Type
}

// Encode serializes a value of type T to the given transport.Writer.
func (c Codec[T, V]) Encode(v T, w transport.Writer) (n int, err error) {
	t := reflect.TypeOf(v)
	dtm, pst := c.typeMap[t]
	if !pst {
		err = NewUnrecognizedType(t)
		return
	}
	n, err = dts.DTMSer.Marshal(dtm, w)
	if err != nil {
		err = NewFailedToMarshalDTM(err)
		return
	}
	bs, err := json.Marshal(v)
	if err != nil {
		err = NewFailedToMarshalJSON(err)
		return
	}
	n1, err := ord.ByteSlice.Marshal(bs, w)
	n += n1
	if err != nil {
		err = NewFailedToMarshalByteSlice(err)
	}
	return
}

// Decode reads a value of type V from the given transport.Reader.
func (c Codec[T, V]) Decode(r transport.Reader) (v V, n int, err error) {
	dtm, n, err := dts.DTMSer.Unmarshal(r)
	if err != nil {
		err = NewFailedToUnmarshalDTM(err)
		return
	}
	if dtm < 0 || dtm >= com.DTM(len(c.dtmSl)) {
		err = NewUnrecognizedDTM(dtm)
		return
	}
	t := c.dtmSl[dtm]
	bs, n1, err := ord.ByteSlice.Unmarshal(r)
	n += n1
	if err != nil {
		err = NewFailedToUnmarshalByteSlice(err)
		return
	}
	ptr := reflect.New(t)
	err = json.Unmarshal(bs, ptr.Interface())
	if err != nil {
		err = NewFailedToUnmarshalJSON(err)
		return
	}
	v = ptr.Elem().Interface().(V)
	return
}
