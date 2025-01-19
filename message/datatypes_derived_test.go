package message

import (
	"net"
	"testing"
)

func TestAddress(t *testing.T) {
	t.Run("SetData_ValidIPv4", func(t *testing.T) {
		ipv4 := net.ParseIP("192.168.1.1")
		address := &Address{}
		err := address.SetData(ipv4)
		if err != nil {
			t.Fatalf("SetData failed: %v", err)
		}
		if !address.isIPv4 {
			t.Errorf("Expected isIPv4 to be true")
		}
		if !address.Data.Equal(ipv4) {
			t.Errorf("Expected Data to be %v, got %v", ipv4, address.Data)
		}
	})

	t.Run("SetData_ValidIPv6", func(t *testing.T) {
		ipv6 := net.ParseIP("2001:db8::ff00:42:8329")
		address := &Address{}
		err := address.SetData(ipv6)
		if err != nil {
			t.Fatalf("SetData failed: %v", err)
		}
		if address.isIPv4 {
			t.Errorf("Expected isIPv4 to be false")
		}
		if !address.Data.Equal(ipv6) {
			t.Errorf("Expected Data to be %v, got %v", ipv6, address.Data)
		}
	})

	t.Run("SetData_InvalidType", func(t *testing.T) {
		address := &Address{}
		err := address.SetData("invalid")
		if err == nil {
			t.Fatalf("Expected SetData to fail with invalid type, but got no error")
		}
	})

	t.Run("Encode_IPv4", func(t *testing.T) {
		ipv4 := net.ParseIP("192.168.1.1")
		address := &Address{}
		err := address.SetData(ipv4)
		if err != nil {
			t.Fatalf("SetData failed: %v", err)
		}
		encoded, err := address.Encode()
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		expected := []byte{0, 1, 192, 168, 1, 1}
		if string(encoded[:6]) != string(expected) {
			t.Errorf("Expected encoded data %v, got %v", expected, encoded[:6])
		}
	})

	t.Run("Encode_IPv6", func(t *testing.T) {
		ipv6 := net.ParseIP("2001:db8::ff00:42:8329")
		address := &Address{}
		err := address.SetData(ipv6)
		if err != nil {
			t.Fatalf("SetData failed: %v", err)
		}
		encoded, err := address.Encode()
		if err != nil {
			t.Fatalf("Encode failed: %v", err)
		}
		expectedPrefix := []byte{0, 2}
		if string(encoded[:2]) != string(expectedPrefix) {
			t.Errorf("Expected encoded prefix %v, got %v", expectedPrefix, encoded[:2])
		}
	})

	t.Run("Decode_ValidIPv4", func(t *testing.T) {
		data := []byte{0, 1, 192, 168, 1, 1}
		address := &Address{}
		err := address.Decode(data)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if !address.isIPv4 {
			t.Errorf("Expected isIPv4 to be true")
		}
		expectedIP := net.IP{192, 168, 1, 1}
		if !address.Data.Equal(expectedIP) {
			t.Errorf("Expected Data to be %v, got %v", expectedIP, address.Data)
		}
	})

	t.Run("Decode_ValidIPv6", func(t *testing.T) {
		data := append([]byte{0, 2}, net.ParseIP("2001:db8::ff00:42:8329").To16()...)
		address := &Address{}
		err := address.Decode(data)
		if err != nil {
			t.Fatalf("Decode failed: %v", err)
		}
		if address.isIPv4 {
			t.Errorf("Expected isIPv4 to be false")
		}
		expectedIP := net.ParseIP("2001:db8::ff00:42:8329")
		if !address.Data.Equal(expectedIP) {
			t.Errorf("Expected Data to be %v, got %v", expectedIP, address.Data)
		}
	})

	t.Run("Decode_InvalidLength", func(t *testing.T) {
		data := []byte{0, 1, 192, 168}
		address := &Address{}
		err := address.Decode(data)
		if err == nil {
			t.Fatalf("Expected Decode to fail with invalid length, but got no error")
		}
	})

	t.Run("String", func(t *testing.T) {
		ipv4 := net.ParseIP("192.168.1.1")
		address := &Address{}
		err := address.SetData(ipv4)
		if err != nil {
			t.Fatalf("SetData failed: %v", err)
		}
		if address.String() != ipv4.String() {
			t.Errorf("Expected String to be %v, got %v", ipv4.String(), address.String())
		}
	})
}
