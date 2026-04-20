package codec

import (
	"bytes"
	"errors"
	"testing"
)

func TestDecodeTruncatedByte(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader(nil)).ReadByte()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedBool(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader(nil)).ReadBool()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedUint16(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader([]byte{0x01})).ReadUint16()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedUint32(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader([]byte{0x01, 0x02})).ReadUint32()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedUint64(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader([]byte{0x01, 0x02, 0x03, 0x04})).ReadUint64()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedStringHeader(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader([]byte{0x01})).ReadString()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedStringPayload(t *testing.T) {
	data := []byte{0x05, 0x00, 0x00, 0x00, 'A', 'B'}
	_, err := NewDecoder(bytes.NewReader(data)).ReadString()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedStringsHeader(t *testing.T) {
	_, err := NewDecoder(bytes.NewReader([]byte{0x01})).ReadStrings()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeTruncatedStringsElement(t *testing.T) {
	data := []byte{0x01, 0x00, 0x00, 0x00}
	_, err := NewDecoder(bytes.NewReader(data)).ReadStrings()
	if !errors.Is(err, ErrDecode) {
		t.Errorf("got %v, want ErrDecode", err)
	}
}

func TestDecodeFromEmpty(t *testing.T) {
	empty := bytes.NewReader(nil)
	if _, err := NewDecoder(empty).ReadUint32(); !errors.Is(err, ErrDecode) {
		t.Errorf("ReadUint32 from empty: got %v, want ErrDecode", err)
	}
}
