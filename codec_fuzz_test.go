package codec_test

import (
	"testing"

	cdctest "github.com/cmd-stream/codec-go/test"
	cdcjson "github.com/cmd-stream/codec-json-go"
)

func FuzzClientCodec_Decode(f *testing.F) {
	reg := cdcjson.NewRegistry(
		cdcjson.WithCmd[any, cdctest.Cmd1](),
		cdcjson.WithCmd[any, cdctest.Cmd2](),
		cdcjson.WithResult[any, cdctest.Result1](),
		cdcjson.WithResult[any, cdctest.Result2](),
	)

	// Seed with some valid data if possible, or just random bytes.
	f.Add([]byte{0, 13, '{', '"', 'A', '"', ':', '1', '0', ',', '"', 'B', '"', ':', '2', '0', '}'}, 0)
	f.Add([]byte{1, 7, '"', 'h', 'e', 'l', 'l', 'o', '"'}, 100)
	f.Add([]byte{255, 0}, 10) // Invalid DTM

	f.Fuzz(func(t *testing.T, data []byte, maxLen int) {
		if maxLen <= 0 || maxLen > 10*1024*1024 {
			return
		}
		c := cdcjson.NewClientCodecWith(reg, cdcjson.WithMaxLen(maxLen))
		cdctest.FuzzDecode(c, data)
	})
}

func FuzzServerCodec_Decode(f *testing.F) {
	reg := cdcjson.NewRegistry(
		cdcjson.WithCmd[any, cdctest.Cmd1](),
		cdcjson.WithCmd[any, cdctest.Cmd2](),
		cdcjson.WithResult[any, cdctest.Result1](),
		cdcjson.WithResult[any, cdctest.Result2](),
	)

	f.Add([]byte{0, 13, '{', '"', 'A', '"', ':', '1', '0', ',', '"', 'B', '"', ':', '2', '0', '}'}, 0)
	f.Add([]byte{1, 7, '"', 'h', 'e', 'l', 'l', 'o', '"'}, 100)

	f.Fuzz(func(t *testing.T, data []byte, maxLen int) {
		if maxLen <= 0 || maxLen > 10*1024*1024 {
			return
		}
		s := cdcjson.NewServerCodecWith(reg, cdcjson.WithMaxLen(maxLen))
		cdctest.FuzzDecode(s, data)
	})
}

func FuzzRoundTrip_Cmd(f *testing.F) {
	var (
		reg = cdcjson.NewRegistry(
			cdcjson.WithCmd[any, cdctest.Cmd1](),
			cdcjson.WithResult[any, cdctest.Result1](),
		)
		client = cdcjson.NewClientCodecWith(reg)
		server = cdcjson.NewServerCodecWith(reg)
	)

	f.Add(10)
	f.Fuzz(func(t *testing.T, x int) {
		cmd := cdctest.Cmd1{A: x, B: x}
		cdctest.VerifyRoundTripCmd(t, client, server, cmd)
	})
}

func FuzzRoundTrip_Result(f *testing.F) {
	var (
		reg = cdcjson.NewRegistry(
			cdcjson.WithCmd[any, cdctest.Cmd1](),
			cdcjson.WithResult[any, cdctest.Result1](),
		)
		client = cdcjson.NewClientCodecWith(reg)
		server = cdcjson.NewServerCodecWith(reg)
	)

	f.Add(10)
	f.Fuzz(func(t *testing.T, x int) {
		res := cdctest.Result1(x)
		cdctest.VerifyRoundTripResult(t, client, server, res)
	})
}
