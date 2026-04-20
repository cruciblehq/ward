package codec

import (
	"io"

	"github.com/cruciblehq/crex"
)

// Sequential binary decoder that reads typed values from an io.Reader.
//
// Each method reads one typed value and advances the read position. Fields
// must be read in the exact order they were written by the Encoder. Methods
// return ErrDecode when the underlying data is truncated or malformed. On
// error the decoder may be left in a partially read state and should not be
// used.
type Decoder struct {
	r io.Reader
}

// Returns a decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Reads a single byte.
func (d *Decoder) ReadByte() (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return 0, crex.Wrap(ErrDecode, err)
	}
	return b[0], nil
}

// Reads a boolean written as a single byte.
func (d *Decoder) ReadBool() (bool, error) {
	var b [1]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return false, crex.Wrap(ErrDecode, err)
	}
	return b[0] != 0, nil
}

// Reads a little-endian uint16.
func (d *Decoder) ReadUint16() (uint16, error) {
	var b [2]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return 0, crex.Wrap(ErrDecode, err)
	}
	return bo.Uint16(b[:]), nil
}

// Reads a little-endian uint32.
func (d *Decoder) ReadUint32() (uint32, error) {
	var b [4]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return 0, crex.Wrap(ErrDecode, err)
	}
	return bo.Uint32(b[:]), nil
}

// Reads a little-endian uint64.
func (d *Decoder) ReadUint64() (uint64, error) {
	var b [8]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return 0, crex.Wrap(ErrDecode, err)
	}
	return bo.Uint64(b[:]), nil
}

// Reads a string written by Encoder.WriteString.
//
// The encoding is a uint32 byte count followed by the raw bytes. An empty
// string is represented by a zero count with no payload. Returns the decoded
// string or ErrDecode if the count header is truncated or the payload is
// shorter than the declared length.
func (d *Decoder) ReadString() (string, error) {
	var b [4]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return "", crex.Wrap(ErrDecode, err)
	}
	n := bo.Uint32(b[:])
	buf := make([]byte, n)
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return "", crex.Wrap(ErrDecode, err)
	}
	return string(buf), nil
}

// Reads a string slice.
//
// Each string is read using ReadString, and the slice is prefixed with a
// uint32 count. A nil or empty slice is represented by a zero count with
// no payload. Returns the decoded slice or ErrDecode if the count header
// or any element is truncated.
func (d *Decoder) ReadStrings() ([]string, error) {
	var b [4]byte
	if _, err := io.ReadFull(d.r, b[:]); err != nil {
		return nil, crex.Wrap(ErrDecode, err)
	}
	n := bo.Uint32(b[:])
	ss := make([]string, n)
	for i := range ss {
		s, err := d.ReadString()
		if err != nil {
			return nil, err
		}
		ss[i] = s
	}
	return ss, nil
}
