package message

import (
	"bytes"
	"testing"
)

func TestEncodeDecodeUint32(t *testing.T) {
	var data uint32 = 123456
	encoded := encode32(data)
	decoded := decode32(encoded, uint32(0))
	if data != decoded {
		t.Fatalf("Expected %d, got %d", data, decoded)
	}
}

func TestOctetString(t *testing.T) {
	t.Run("SetData", func(t *testing.T) {
		tests := []struct {
			name    string
			input   interface{}
			want    []byte
			wantErr bool
		}{
			{"SetData with []byte", []byte("test"), []byte("test"), false},
			{"SetData with string", "test", []byte("test"), false},
			{"SetData with invalid type", 123, nil, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := &OctetString{}
				err := o.SetData(tt.input)

				if (err != nil) != tt.wantErr {
					t.Errorf("SetData() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr && !bytes.Equal(o.Data, tt.want) {
					t.Errorf("SetData() got = %v, want %v", o.Data, tt.want)
				}
			})
		}
	})

	t.Run("Length", func(t *testing.T) {
		o := &OctetString{Data: []byte("test")}
		got := o.Length()
		want := uint32(4)

		if got != want {
			t.Errorf("Length() got = %v, want %v", got, want)
		}
	})

	t.Run("Encode", func(t *testing.T) {
		// getPadding = func(length int) int { return (4 - (length % 4)) % 4 } // Mock padding calculation

		tests := []struct {
			name       string
			data       []byte
			minLength  uint32
			wantLength int
		}{
			{"Encode with padding", []byte("test"), 8, 8},
			{"Encode with minLength", []byte("test"), 12, 12},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				o := &OctetString{
					Data:       tt.data,
					min_length: tt.minLength,
				}
				encoded, err := o.Encode()
				if err != nil {
					t.Errorf("Encode() error = %v", err)
					return
				}

				if len(encoded) != tt.wantLength {
					t.Errorf("Encode() got length = %v, want %v", len(encoded), tt.wantLength)
				}

				if !bytes.HasPrefix(encoded, tt.data) {
					t.Errorf("Encode() encoded data does not match input data")
				}
			})
		}
	})

	t.Run("Decode", func(t *testing.T) {
		data := []byte("decoded data")
		o := &OctetString{}
		err := o.Decode(data)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
		}

		if !bytes.Equal(o.Data, data) {
			t.Errorf("Decode() got = %v, want %v", o.Data, data)
		}
	})

	t.Run("String", func(t *testing.T) {
		o := &OctetString{Data: []byte("test string")}
		got := o.String()
		want := "test string"

		if got != want {
			t.Errorf("String() got = %v, want %v", got, want)
		}
	})

	t.Run("Type", func(t *testing.T) {
		o := &OctetString{}
		want := OctetStringType
		got := o.Type()

		if got != want {
			t.Errorf("Type() got = %v, want %v", got, want)
		}
	})
}
