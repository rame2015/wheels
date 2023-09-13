package wheels

import "errors"

var (
	ErrServiceAlreadyExists   = errors.New("service already exists")
	ErrUnknownService         = errors.New("unknown service")
	ErrServiceNotImplementsAs = errors.New("service not implements as")
	ErrInvalidAsType          = errors.New("invalid as type")
	ErrInvalidCtorType        = errors.New("invalid ctor type")
	ErrInvalidZeroType        = errors.New("invalid zero type")
	ErrInvalidInvokeType      = errors.New("invalid invoke type")
)
