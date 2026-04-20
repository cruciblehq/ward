package codec

import (
	"encoding/binary"
)

// Byte order for multi-byte values in binary-encoded formats.
//
// Using little-endian byte order allows the encoder and decoder to read and
// write multi-byte values directly from memory without needing to rearrange
// bytes, which is more efficient on x86-64 and arm64 targets. Additionally,
// little-endian is the byte order used by BPF maps, which may be relevant
// for certain subsystems that interact with BPF.
var bo = binary.LittleEndian

// Paired encoding and decoding of a single type.
//
// Subsystems implement this interface so their state types can be serialized
// and deserialized through a sequential binary stream. EncodeTo writes fields
// in order; DecodeFrom reads them back in the same order.
type Codec interface {

	// Writes all fields to the encoder in declaration order.
	//
	// The implementation must write exactly the same sequence of typed values
	// that DecodeFrom expects to read, so the two methods form a matched pair.
	EncodeTo(e *Encoder) error

	// Reads all fields from the decoder in declaration order.
	//
	// Returns ErrDecode if any field is truncated or malformed. On error the
	// receiver may be left in a partially populated state and should not be
	// used.
	DecodeFrom(d *Decoder) error
}
