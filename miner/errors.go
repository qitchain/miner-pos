package miner

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

var (
	ErrUnauthenticated  = status.Error(codes.Unauthenticated, "unauthenticated")
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
	ErrBadRequest       = status.Error(codes.InvalidArgument, "bad request")
	ErrNotFoundData     = status.Error(codes.NotFound, "not found data")
	ErrBadMiningInfo    = errors.New("bad mining info")
)

// RawRpcError raw rpc error
type RawRpcError struct {
	code    int
	message string
}

func (err RawRpcError) Code() int {
	return err.code
}

func (err RawRpcError) Message() string {
	return err.message
}

func (err RawRpcError) Error() string {
	return err.Message() + "(" + strconv.Itoa(err.Code()) + ")"
}

// WrapRawRpcError wrap raw rpc message to error
func WrapRawRpcError(code int, message string) error {
	return RawRpcError{
		code:    code,
		message: message,
	}
}
