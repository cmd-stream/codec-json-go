package codec

import "encoding/json"

// Serializer implements the gnrc.Serializer interface using JSON.
type Serializer[T, V any] struct{}

// Marshal encodes the given value of type T into a JSON byte slice.
func (s Serializer[T, V]) Marshal(t T) (bs []byte, err error) {
	return json.Marshal(t)
}

// Unmarshal decodes the JSON byte slice into the given value of type V.
func (s Serializer[T, V]) Unmarshal(bs []byte, v V) (err error) {
	return json.Unmarshal(bs, v)
}
