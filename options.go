package codec

import cdc "github.com/cmd-stream/codec-go"

// SetOption is a function that sets an option for the codec.
type SetOption = cdc.SetOption

// WithMaxLen returns a SetOption that sets the maximum length of the encoded
// message byte slice.
func WithMaxLen(maxLen int) SetOption {
	return cdc.WithMaxLen(maxLen)
}
