package codec

import (
	"github.com/cmd-stream/cmd-stream-go/core"
	cdc "github.com/cmd-stream/codec-go"
)

// Registry is a collection of Command and Result types.
//
// The order of registration matters — two registries created with the same
// types in a different order will produce codecs that are not compatible.
type Registry[T any] = cdc.Registry[T]

// NewRegistry creates a new Registry with the provided options.
func NewRegistry[T any](opts ...func(*Registry[T])) *Registry[T] {
	return cdc.NewRegistry(opts...)
}

// WithCmd returns an option that registers a command type C in the Registry.
//
// The order in which WithCmd options are passed to NewRegistry matters.
func WithCmd[T any, C core.Cmd[T]]() func(*Registry[T]) {
	return cdc.WithCmd[T, C]()
}

// WithResult returns an option that registers a result type R in the Registry.
//
// The order in which WithResult options are passed to NewRegistry matters.
func WithResult[T any, R core.Result]() func(*Registry[T]) {
	return cdc.WithResult[T, R]()
}
