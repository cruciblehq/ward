// Package codec provides low-level helpers for encoding and decoding
// sequential binary data.
//
// Subsystem state types use the Encoder and Decoder to serialize their
// fields in a portable little-endian format. The Codec interface ties
// the two together so that any type can describe its own wire layout.
//
// All multi-byte values use little-endian byte order. Strings are
// written as a uint32 byte count followed by the raw bytes, with no
// trailing NUL. The encoder and decoder are strictly sequential: fields
// must be read back in the same order they were written. Neither type
// is safe for concurrent use.
//
// Format-level concerns (headers, slot tables, file I/O) live outside
// this package.
package codec
