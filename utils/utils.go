package utils

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode(data []byte) error
}

func Encode(e Encoder) ([]byte, error) {
	return e.Encode()
}

func Decode(d Decoder, data []byte) error {
	return d.Decode(data)
}

// ToBytes converts a uint32 to a byte slice.
func ToBytes(value uint32, count int) []byte {
	result := make([]byte, count)
	for i := 0; i < count; i++ {
		sh := (count - i - 1) * 8
		result[i] = byte(value >> uint(sh))
	}
	return result
}

// FromBytes converts a byte slice to a uint32.
func FromBytes(data []byte) uint32 {
	var result uint32
	for i := 0; i < len(data); i++ {
		result = result<<8 | uint32(data[i])
	}
	return result
}
