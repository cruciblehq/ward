package codec

import (
	"io"

	"github.com/cruciblehq/crex"
)

// Sequential binary encoder backed by an io.Writer.
//
// Each method appends one typed value and returns any write error. The caller
// is responsible for writing fields in a fixed order so that the Decoder can
// read them back.
type Encoder struct {
	w io.Writer // Underlying writer for encoded bytes.
}

// Returns an encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Writes a string.
//
// The encoding is a uint32 byte count followed by the raw bytes. An empty
// string is written as a zero count with no payload.
func (e *Encoder) WriteString(s string) error {
	var b [4]byte
	bo.PutUint32(b[:], uint32(len(s)))
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	if _, err := io.WriteString(e.w, s); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	return nil
}

// Writes a string slice.
//
// The encoding is a uint32 count followed by each string in length-prefixed
// form. A nil or empty slice is written as a zero count with no payload.
func (e *Encoder) WriteStrings(ss []string) error {
	var b [4]byte
	bo.PutUint32(b[:], uint32(len(ss)))
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	for _, s := range ss {
		if err := e.WriteString(s); err != nil {
			return err
		}
	}
	return nil
}

// Writes a little-endian uint32.
func (e *Encoder) WriteUint32(v uint32) error {
	var b [4]byte
	bo.PutUint32(b[:], v)
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	return nil
}

// Writes a little-endian uint64.
func (e *Encoder) WriteUint64(v uint64) error {
	var b [8]byte
	bo.PutUint64(b[:], v)
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	return nil
}

// Writes a boolean as a single byte (0 or 1).
func (e *Encoder) WriteBool(v bool) error {
	var b [1]byte
	if v {
		b[0] = 1
	}
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	return nil
}

// Writes a single byte.
func (e *Encoder) WriteByte(v byte) error {
	b := [1]byte{v}
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	return nil
}

// Writes a little-endian uint16.
func (e *Encoder) WriteUint16(v uint16) error {
	var b [2]byte
	bo.PutUint16(b[:], v)
	if _, err := e.w.Write(b[:]); err != nil {
		return crex.Wrap(ErrEncode, err)
	}
	return nil
}
