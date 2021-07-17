package libravatar

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type ImageSize uint16

const (
	SizeMinimum ImageSize = 1
	SizeMaximum ImageSize = 512
	SizeDefault ImageSize = 80
)
