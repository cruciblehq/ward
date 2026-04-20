package cli

import "errors"

var (
	ErrParse  = errors.New("failed to parse affordance")
	ErrBuild  = errors.New("build failed")
	ErrOutput = errors.New("failed to write output")
)
