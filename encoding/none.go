package encoding

import (
	"fmt"
)

// NoCodec encodes/decodes Go values to/from []byte.
// You can use encoding.NoCodec instead of creating an instance of this struct.
type NoCodec struct{}

// Marshal encodes a Go value
func (c NoCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("only byte slices supported")
	}
	return b, nil
}

// Unmarshal decodes a gob value
func (c NoCodec) Unmarshal(data []byte, v interface{}) error {
	b, ok := v.(*[]byte)
	if !ok {
		return fmt.Errorf("only byte slices supported")
	}
	*b = data
	return nil
}
