// Package codecjson provides a JSON-based codec implementation for
// cmd-stream-go.
package codecjson

import (
	"encoding/json"
	"reflect"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/dts-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
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

// NewCodec constructs a Codec with explicit type mappings.
//
// The first slice lists types that can be encoded,
// and the second slice lists types that can be decoded.
func NewCodec[T, V any](types1 []reflect.Type, types2 []reflect.Type) (
	codec Codec[T, V],
) {
	codec = Codec[T, V]{
		typeMap: make(map[reflect.Type]com.DTM),
		dtmMap:  make(map[com.DTM]reflect.Type),
	}
	for i, t := range types1 {
		codec.typeMap[t] = com.DTM(i)
	}
	for i, t := range types2 {
		codec.dtmMap[com.DTM(i)] = t
	}
	return
}

// Codec provides JSON serialization and deserialization for registered types.
//
// T represents the types that can be encoded, and V represents the types that
// can be decoded.
type Codec[T, V any] struct {
	typeMap map[reflect.Type]com.DTM
	dtmMap  map[com.DTM]reflect.Type
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
	t, pst := c.dtmMap[dtm]
	if !pst {
		err = NewUnrecognizedDTM(dtm)
		return
	}
	bs, n1, err := ord.ByteSlice.Unmarshal(r)
	n += n1
	if err != nil {
		err = NewFailedToUnmarshalByteSlice(err)
		return
	}
	val := reflect.New(t).Interface()
	err = json.Unmarshal(bs, val)
	if err != nil {
		err = NewFailedToUnmarshalJSON(err)
		return
	}
	v = reflect.ValueOf(val).Elem().Interface().(V)
	return
}
