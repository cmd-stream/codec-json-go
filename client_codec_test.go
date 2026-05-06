package codec_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	cmock "github.com/cmd-stream/cmd-stream-go/test/mock"
	cdcjson "github.com/cmd-stream/codec-json-go"
	cdctest "github.com/cmd-stream/codec-go/test"
	com "github.com/mus-format/common-go"
	assertfatal "github.com/ymz-ncnk/assert/fatal"
)

func TestClientCodec_Encoding(t *testing.T) {
	var (
		wantDTM   = 0
		cmd       = cdctest.Cmd1{A: 10, B: 20}
		wantBs, _ = json.Marshal(cmd)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		writer    = cmock.NewWriter()
	)
	writer.RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(wantDTM))
			return nil
		},
	).RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(wantLen))
			return nil
		},
	).RegisterWrite(
		func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(t, p, wantBs)
			return len(p), nil
		},
	)
	codec := cdcjson.NewClientCodec[any](
		[]reflect.Type{
			reflect.TypeFor[cdctest.Cmd1](),
			reflect.TypeFor[cdctest.Cmd2](),
		},
		[]reflect.Type{
			reflect.TypeFor[cdctest.Result1](),
			reflect.TypeFor[cdctest.Result2](),
		},
	)
	n, err := codec.Encode(cmd, writer)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
}

func TestClientCodec_EncodeError(t *testing.T) {
	var (
		cmd     = cdctest.Cmd1{A: 10, B: 20}
		wantErr = errors.New("write error")
		writer  = cmock.NewWriter()
	)
	writer.RegisterWriteByte(func(b byte) error {
		return wantErr
	})
	codec := cdcjson.NewClientCodec[any](
		[]reflect.Type{reflect.TypeFor[cdctest.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[cdctest.Result1]()},
	)
	_, err := codec.Encode(cmd, writer)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdcjson.ErrorPrefix)], cdcjson.ErrorPrefix)
}

func TestClientCodec_Decoding(t *testing.T) {
	var (
		wantDTM   = 1
		wantV     = cdctest.Result2{Y: "hello"}
		wantBs, _ = json.Marshal(wantV)
		wantLen   = len(wantBs)
		wantN     = 1 + 1 + wantLen
		reader    = cmock.NewReader()
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return byte(wantDTM), nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(wantLen), nil },
	).RegisterRead(
		func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return wantLen, nil
		},
	)
	codec := cdcjson.NewClientCodec[any](
		[]reflect.Type{
			reflect.TypeFor[cdctest.Cmd1](),
			reflect.TypeFor[cdctest.Cmd2](),
		},
		[]reflect.Type{
			reflect.TypeFor[cdctest.Result1](),
			reflect.TypeFor[cdctest.Result2](),
		},
	)
	v, n, err := codec.Decode(reader)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
	assertfatal.EqualDeep(t, v, wantV)
}

func TestClientCodec_DecodeError(t *testing.T) {
	var (
		wantErr = errors.New("read error")
		reader  = cmock.NewReader()
	)
	reader.RegisterReadByte(func() (b byte, err error) {
		return 0, wantErr
	})
	codec := cdcjson.NewClientCodec[any](
		[]reflect.Type{reflect.TypeFor[cdctest.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[cdctest.Result1]()},
	)
	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdcjson.ErrorPrefix)], cdcjson.ErrorPrefix)
}

func TestClientCodecWith(t *testing.T) {
	var (
		wantDTM   = 0
		cmd       = cdctest.Cmd1{A: 10, B: 20}
		wantBs, _ = json.Marshal(cmd)
		writer    = cmock.NewWriter()

		wantResultDTM   = 0
		wantV           = cdctest.Result1(10)
		wantResultBs, _ = json.Marshal(wantV)
		reader          = cmock.NewReader()
	)
	writer.RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(wantDTM))
			return nil
		},
	).RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(len(wantBs)))
			return nil
		},
	).RegisterWrite(
		func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(t, p, wantBs)
			return len(p), nil
		},
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return byte(wantResultDTM), nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(len(wantResultBs)), nil },
	).RegisterRead(
		func(p []byte) (n int, err error) {
			copy(p, wantResultBs)
			return len(wantResultBs), nil
		},
	)

	registry := cdcjson.NewRegistry(
		cdcjson.WithCmd[any, cdctest.Cmd1](),
		cdcjson.WithResult[any, cdctest.Result1](),
	)
	codec := cdcjson.NewClientCodecWith(registry)

	// Verify Encoding
	_, err := codec.Encode(cmd, writer)
	assertfatal.EqualError(t, err, nil)

	// Verify Decoding
	v, _, err := codec.Decode(reader)
	assertfatal.EqualError(t, err, nil)
	assertfatal.EqualDeep(t, v, wantV)
}

func TestClientCodec_MaxLenOption(t *testing.T) {
	var (
		maxLen = 5
		reader = cmock.NewReader()
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return 0, nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(10), nil },
	)

	codec := cdcjson.NewClientCodec[any](
		[]reflect.Type{reflect.TypeFor[cdctest.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[cdctest.Result2]()},
		cdcjson.WithMaxLen(maxLen),
	)

	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, com.ErrTooLargeLength), true)
}
