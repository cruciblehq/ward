package codec

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestByteWireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteByte(0x42); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x42}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

func TestBoolWireFormat(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.WriteBool(false); err != nil {
		t.Fatal(err)
	}
	if err := enc.WriteBool(true); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00, 0x01}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

func TestUint16WireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteUint16(0x0201); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x01, 0x02}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

func TestUint32WireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteUint32(0x04030201); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x01, 0x02, 0x03, 0x04}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

func TestUint64WireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteUint64(0x0807060504030201); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

func TestStringWireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteString("AB"); err != nil {
		t.Fatal(err)
	}
	want := []byte{0x02, 0x00, 0x00, 0x00, 'A', 'B'}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

func TestStringsWireFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := NewEncoder(&buf).WriteStrings([]string{"A", "BC"}); err != nil {
		t.Fatal(err)
	}
	want := []byte{
		0x02, 0x00, 0x00, 0x00, // count = 2
		0x01, 0x00, 0x00, 0x00, 'A', // "A"
		0x02, 0x00, 0x00, 0x00, 'B', 'C', // "BC"
	}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("wire bytes = %x, want %x", buf.Bytes(), want)
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

type limitWriter struct {
	limit   int
	written int
}

func (w *limitWriter) Write(p []byte) (int, error) {
	remaining := w.limit - w.written
	if remaining <= 0 {
		return 0, io.ErrShortWrite
	}
	if len(p) > remaining {
		w.written += remaining
		return remaining, io.ErrShortWrite
	}
	w.written += len(p)
	return len(p), nil
}

func TestEncodeWriteErrorUint16(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteUint16(1)
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestEncodeWriteErrorUint32(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteUint32(1)
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestEncodeWriteErrorUint64(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteUint64(1)
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestEncodeWriteErrorBool(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteBool(true)
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestEncodeWriteErrorByte(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteByte(0x42)
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestEncodeWriteErrorString(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteString("hello")
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestEncodeWriteErrorStrings(t *testing.T) {
	err := NewEncoder(errWriter{}).WriteStrings([]string{"a"})
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}

func TestStringWritePayloadError(t *testing.T) {
	w := &limitWriter{limit: 4}
	err := NewEncoder(w).WriteString("hello")
	if !errors.Is(err, ErrEncode) {
		t.Errorf("got %v, want ErrEncode", err)
	}
}
