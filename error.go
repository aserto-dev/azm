package azm

import (
	"net/http"

	cerr "github.com/aserto-dev/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidSchemaVersion = cerr.NewAsertoError("E10000", codes.InvalidArgument, http.StatusBadRequest, "invalid or unsupported schema version")
	ErrObjectNotFound       = cerr.NewAsertoError("E10001", codes.NotFound, http.StatusNotFound, "object type not found")
	ErrRelationNotFound     = cerr.NewAsertoError("E10002", codes.NotFound, http.StatusNotFound, "relation type not found")
	ErrPermissionNotFound   = cerr.NewAsertoError("E10003", codes.NotFound, http.StatusNotFound, "permission type not found")
)
