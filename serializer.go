package codec

import "encoding/json"

type Serializer[T, V any] struct{}

func (s Serializer[T, V]) Marshal(t T) (bs []byte, err error) {
	return json.Marshal(t)
}

func (s Serializer[T, V]) Unmarshal(bs []byte, v V) (err error) {
	return json.Unmarshal(bs, v)
}
