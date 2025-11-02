package codec_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/cmd-stream/codec-json-go"
	tmocks "github.com/cmd-stream/testkit-go/mocks/transport"

	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestServerCodec(t *testing.T) {
	t.Run("Encoding should work", func(t *testing.T) {
		wantDTM := 0
		result := Result1{X: 10}
		wantBs, err := json.Marshal(result)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		c := codec.NewServerCodec[any](
			[]reflect.Type{
				reflect.TypeFor[Cmd1](),
				reflect.TypeFor[Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[Result1](),
				reflect.TypeFor[Result2](),
			},
		)

		w := tmocks.NewWriter().RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(b, byte(wantDTM), t)
			return nil
		}).RegisterWriteByte(func(b byte) error {
			assertfatal.Equal(b, byte(wantLen), t)
			return nil
		}).RegisterWrite(func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(p, wantBs, t)
			return len(p), nil
		})

		n, err := c.Encode(result, w)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
	})

	t.Run("Decoding should work", func(t *testing.T) {
		wantDTM := 1
		wantV := Cmd2{Y: "hello"}
		wantBs, err := json.Marshal(wantV)
		assertfatal.EqualError(err, nil, t)
		wantLen := len(wantBs)
		wantN := 1 + 1 + wantLen

		c := codec.NewServerCodec[any](
			[]reflect.Type{
				reflect.TypeFor[Cmd1](),
				reflect.TypeFor[Cmd2](),
			},
			[]reflect.Type{
				reflect.TypeFor[Result1](),
				reflect.TypeFor[Result2](),
			},
		)

		r := tmocks.NewReader().RegisterReadByte(func() (b byte, err error) {
			return byte(wantDTM), nil
		}).RegisterReadByte(func() (b byte, err error) {
			return byte(wantLen), nil
		}).RegisterRead(func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		})

		v, n, err := c.Decode(r)
		assertfatal.EqualError(err, nil, t)
		assertfatal.Equal(n, wantN, t)
		assertfatal.EqualDeep(v, wantV, t)
	})
}
