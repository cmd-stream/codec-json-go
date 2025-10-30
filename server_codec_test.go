package codec

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cmd-stream/codec-json-go/testdata"

	tmock "github.com/cmd-stream/transport-go/testdata/mock"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestServerCodec(t *testing.T) {
	t.Run("Encoding should work", func(t *testing.T) {
		wantDTM := 0
		result := testdata.Result1{X: 10}
		wantBs, err := json.Marshal(result)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		codec := NewServerCodec[any](
			[]reflect.Type{
				reflect.TypeFor[testdata.Cmd1](),
				reflect.TypeFor[testdata.Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.Result1](),
				reflect.TypeFor[testdata.Result2](),
			},
		)

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

		n, err := codec.Encode(result, w)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
	})

	t.Run("Decoding should work", func(t *testing.T) {
		wantDTM := 1
		wantV := testdata.Cmd2{Y: "hello"}
		wantBs, err := json.Marshal(wantV)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		codec := NewServerCodec[any](
			[]reflect.Type{
				reflect.TypeFor[testdata.Cmd1](),
				reflect.TypeFor[testdata.Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[testdata.Result1](),
				reflect.TypeFor[testdata.Result2](),
			},
		)

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
}
