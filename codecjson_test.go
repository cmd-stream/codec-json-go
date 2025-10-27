package codecjson

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/cmd-stream/codec-json-go/testdata"
	tmock "github.com/cmd-stream/transport-go/testdata/mock"
	com "github.com/mus-format/common-go"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestCodecJSON(t *testing.T) {
	t.Run("Encoding should work", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		wantDTM := 0
		v := testdata.MyStruct1{X: 10}
		wantBs, err := json.Marshal(v)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		w := tmock.NewWriter().RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(b, byte(wantDTM), t)
			return nil
		}).RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(b, byte(wantLen), t)
			return nil
		}).RegisterWrite(func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(p, wantBs, t)
			return len(p), nil
		})

		n, err := codec.Encode(v, w)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
	})

	t.Run("Failed to marshal DTM", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		writeErr := errors.New("failed to write DTM")
		wantErr := NewFailedToMarshalDTM(writeErr)

		w := tmock.NewWriter().RegisterWriteByte(func(b byte) error {
			return writeErr
		})
		n, err := codec.Encode(testdata.MyStruct1{}, w)
		assertfatal.EqualError(err, wantErr, t)
		assertfatal.Equal(n, 0, t)
	})

	t.Run("Decoding should work", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		wantDTM := 1
		wantV := testdata.MyStruct2{Y: "hello"}
		wantBs, err := json.Marshal(wantV)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		r := tmock.NewReader().RegisterReadByte(func() (b byte, err error) {
			return byte(wantDTM), nil
		}).RegisterReadByte(func() (b byte, err error) {
			return byte(wantLen), nil
		}).RegisterRead(func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		})

		v, n, err := codec.Decode(r)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
		assertfatal.EqualDeep(v, wantV, t)
	})

	t.Run("Unrecognized type", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		v := testdata.MyStruct3{Z: 3.14}
		wantType := reflect.TypeOf(v)
		wantErr := NewUnrecognizedType(wantType)

		w := tmock.NewWriter()

		n, err := codec.Encode(v, w)
		assertfatal.EqualError(err, wantErr, t)
		assertfatal.Equal(n, 0, t)
	})

	t.Run("Failed to marshal byte slice", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		writeErr := errors.New("failed to write byte slice length")
		wantErr := NewFailedToMarshalByteSlice(writeErr)

		w := tmock.NewWriter().RegisterWriteByte(func(b byte) error {
			return nil
		}).RegisterWriteByte(func(b byte) error {
			return writeErr
		})

		n, err := codec.Encode(testdata.MyStruct1{}, w)
		assertfatal.EqualError(err, wantErr, t)
		assertfatal.Equal(n, 1, t)
	})

	t.Run("Failed to unmarshal DTM", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		readErr := errors.New("failed to read DTM")
		wantErr := NewFailedToUnmarshalDTM(readErr)

		r := tmock.NewReader().RegisterReadByte(func() (b byte, err error) {
			return 0, readErr
		})

		_, n, err := codec.Decode(r)
		assertfatal.EqualError(err, wantErr, t)
		assertfatal.Equal(n, 0, t)
	})

	t.Run("Unrecognized DTM", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		const unrecognizedDTM com.DTM = 99
		wantErr := NewUnrecognizedDTM(unrecognizedDTM)

		r := tmock.NewReader().RegisterReadByte(func() (b byte, err error) {
			return byte(unrecognizedDTM), nil
		})

		_, n, err := codec.Decode(r)
		assertfatal.EqualError(err, wantErr, t)
		assertfatal.Equal(n, 1, t)
	})

	t.Run("Failed to unmarshal byte slice", func(t *testing.T) {
		codec := NewCodec[testdata.MyInterface, testdata.MyInterface](
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.MyStruct1](),
				reflect.TypeFor[testdata.MyStruct2](),
			},
		)

		readErr := errors.New("failed to read byte slice")
		wantErr := NewFailedToUnmarshalByteSlice(readErr)

		r := tmock.NewReader().RegisterReadByte(func() (b byte, err error) {
			return 0, nil
		}).RegisterReadByte(func() (b byte, err error) {
			return 0, readErr
		})

		_, n, err := codec.Decode(r)
		assertfatal.EqualError(err, wantErr, t)
		assertfatal.Equal(n, 1, t)
	})
}
