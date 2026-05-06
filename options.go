package codec

import cdc "github.com/cmd-stream/codec-go"

// SetOption is a function that sets an option for the codec.
type SetOption = cdc.SetOption

// WithMaxLen sets the maximum length of the byte slice.
func WithMaxLen(maxLen int) SetOption {
	return cdc.WithMaxLen(maxLen)
}
