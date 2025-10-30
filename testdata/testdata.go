package testdata

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
)

type Cmd1 struct {
	X int
}

func (c Cmd1) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver any, proxy core.Proxy,
) error {
	return nil
}

type Cmd2 struct {
	Y string
}

func (c Cmd2) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver any, proxy core.Proxy,
) error {
	return nil
}

type Result1 struct {
	X int
}

func (r Result1) LastOne() bool {
	return true
}

type Result2 struct {
	Y string
}

func (r Result2) LastOne() bool {
	return true
}
