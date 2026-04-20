package codec

import "errors"

var (
	ErrEncode = errors.New("encode failed")
	ErrDecode = errors.New("decode failed")
)
