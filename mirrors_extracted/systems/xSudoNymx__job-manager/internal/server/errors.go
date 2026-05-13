package server

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xSudoNymx/job-manager/pkg/manager"
)

// toStatus converts a library error to a gRPC status error.
func toStatus(err error) error {
	if err == nil {
		return nil
	}

	code := toCode(err)
	return status.Error(code, err.Error())
}

// toCode returns the appropriate gRPC code for an error.
func toCode(err error) codes.Code {
	switch {
	case errors.Is(err, manager.ErrJobNotFound):
		return codes.NotFound

	case errors.Is(err, manager.ErrStorageUnavailable),
		errors.Is(err, manager.ErrStartFailed),
		errors.Is(err, manager.ErrStopFailed),
		errors.Is(err, manager.ErrSubscribeFailed),
		errors.Is(err, manager.ErrShutdownFailed):
		return codes.Internal

	default:
		return codes.Internal
	}
}
