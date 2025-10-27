package codecjson

import (
	"fmt"
	"reflect"

	com "github.com/mus-format/common-go"
)

const errorPrefix = "codecjson: "

// NewUnrecognizedType returns an error indicating that an unsupported type
// was encountered during encoding.
func NewUnrecognizedType(t reflect.Type) error {
	return fmt.Errorf(errorPrefix+"unrecognized type: %T", t)
}

// NewFailedToMarshalDTM returns an error indicating that the data type marker
// (DTM) could not be marshaled.
func NewFailedToMarshalDTM(err error) error {
	return fmt.Errorf(errorPrefix+"failed to marshal DTM: %w", err)
}

// NewFailedToMarshalJSON returns an error indicating that JSON marshaling
// failed.
func NewFailedToMarshalJSON(err error) error {
	return fmt.Errorf(errorPrefix+"failed to marshal JSON: %w", err)
}

// NewFailedToMarshalByteSlice returns an error indicating that a byte slice
// could not be marshaled.
func NewFailedToMarshalByteSlice(err error) error {
	return fmt.Errorf(errorPrefix+"failed to marshal byte slice: %w", err)
}

// NewFailedToUnmarshalDTM returns an error indicating that the data type
// marker (DTM) could not be unmarshaled.
func NewFailedToUnmarshalDTM(err error) error {
	return fmt.Errorf(errorPrefix+"failed to unmarshal DTM: %w", err)
}

// NewUnrecognizedDTM returns an error indicating that an unknown data type
// marker (DTM) was received.
func NewUnrecognizedDTM(dtm com.DTM) error {
	return fmt.Errorf(errorPrefix+"unrecognized DTM: %v", dtm)
}

// NewFailedToUnmarshalByteSlice returns an error indicating that a byte slice
// could not be unmarshaled.
func NewFailedToUnmarshalByteSlice(err error) error {
	return fmt.Errorf(errorPrefix+"failed to unmarshal byte slice: %w", err)
}

// NewFailedToUnmarshalJSON returns an error indicating that JSON unmarshaling
// failed.
func NewFailedToUnmarshalJSON(err error) error {
	return fmt.Errorf(errorPrefix+"failed to unmarshal JSON: %w", err)
}
