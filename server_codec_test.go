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

func TestServerCodec_Encoding(t *testing.T) {
	var (
		wantDTM   = 0
		result    = cdctest.Result1(10)
		wantBs, _ = json.Marshal(result)
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
	codec := cdcjson.NewServerCodec[any](
		[]reflect.Type{
			reflect.TypeFor[cdctest.Cmd1](),
			reflect.TypeFor[cdctest.Cmd2](),
		},
		[]reflect.Type{
			reflect.TypeFor[cdctest.Result1](),
			reflect.TypeFor[cdctest.Result2](),
		},
	)
	n, err := codec.Encode(result, writer)
	assertfatal.EqualError(t, err, nil)
	assertfatal.Equal(t, n, wantN)
}

func TestServerCodec_EncodeError(t *testing.T) {
	var (
		result  = cdctest.Result1(10)
		wantErr = errors.New("write error")
		writer  = cmock.NewWriter()
	)
	writer.RegisterWriteByte(func(b byte) error {
		return wantErr
	})
	codec := cdcjson.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[cdctest.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[cdctest.Result1]()},
	)
	_, err := codec.Encode(result, writer)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdcjson.ErrorPrefix)], cdcjson.ErrorPrefix)
}

func TestServerCodec_Decoding(t *testing.T) {
	var (
		wantDTM   = 1
		wantV     = cdctest.Cmd2("hello")
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
	codec := cdcjson.NewServerCodec[any](
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

func TestServerCodec_DecodeError(t *testing.T) {
	var (
		wantErr = errors.New("read error")
		reader  = cmock.NewReader()
	)
	reader.RegisterReadByte(func() (b byte, err error) {
		return 0, wantErr
	})
	codec := cdcjson.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[cdctest.Cmd1]()},
		[]reflect.Type{reflect.TypeFor[cdctest.Result1]()},
	)
	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, wantErr), true)
	assertfatal.EqualDeep(t, err.Error()[:len(cdcjson.ErrorPrefix)], cdcjson.ErrorPrefix)
}

func TestServerCodecWith(t *testing.T) {
	var (
		wantCmdDTM = 0
		cmd        = cdctest.Cmd1{A: 10, B: 20}
		wantBs, _  = json.Marshal(cmd)
		reader     = cmock.NewReader()

		wantResultDTM   = 0
		result          = cdctest.Result1(10)
		wantResultBs, _ = json.Marshal(result)
		writer          = cmock.NewWriter()
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return byte(wantCmdDTM), nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(len(wantBs)), nil },
	).RegisterRead(
		func(p []byte) (n int, err error) {
			copy(p, wantBs)
			return len(wantBs), nil
		},
	)
	writer.RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(wantResultDTM))
			return nil
		},
	).RegisterWriteByte(
		func(b byte) error {
			assertfatal.Equal(t, b, byte(len(wantResultBs)))
			return nil
		},
	).RegisterWrite(
		func(p []byte) (n int, err error) {
			assertfatal.EqualDeep(t, p, wantResultBs)
			return len(p), nil
		},
	)

	registry := cdcjson.NewRegistry(
		cdcjson.WithCmd[any, cdctest.Cmd1](),
		cdcjson.WithResult[any, cdctest.Result1](),
	)
	codec := cdcjson.NewServerCodecWith(registry)

	// Verify Decoding
	v, _, err := codec.Decode(reader)
	assertfatal.EqualError(t, err, nil)
	assertfatal.EqualDeep(t, v, cmd)

	// Verify Encoding
	_, err = codec.Encode(result, writer)
	assertfatal.EqualError(t, err, nil)
}

func TestServerCodec_MaxLenOption(t *testing.T) {
	var (
		maxLen = 5
		reader = cmock.NewReader()
	)
	reader.RegisterReadByte(
		func() (b byte, err error) { return 0, nil },
	).RegisterReadByte(
		func() (b byte, err error) { return byte(10), nil },
	)
	codec := cdcjson.NewServerCodec[any](
		[]reflect.Type{reflect.TypeFor[cdctest.Cmd2]()},
		[]reflect.Type{reflect.TypeFor[cdctest.Result1]()},
		cdcjson.WithMaxLen(maxLen),
	)

	_, _, err := codec.Decode(reader)
	assertfatal.EqualDeep(t, errors.Is(err, com.ErrTooLargeLength), true)
}
