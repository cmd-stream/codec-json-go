package cmds

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
