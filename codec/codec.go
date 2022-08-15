package codec

// Codec declares encode and decode functions
type Codec interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}
