package cmds

import (
	"context"
	"time"

	"github.com/cmd-stream/core-go"
)

type Cmd2 struct {
	Y string
}

func (c Cmd2) Exec(ctx context.Context, seq core.Seq, at time.Time,
	receiver any, proxy core.Proxy,
) error {
	return nil
}
