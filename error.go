package azm

import (
	"net/http"

	cerr "github.com/aserto-dev/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidObjectType      = cerr.NewAsertoError("E21001", codes.InvalidArgument, http.StatusBadRequest, "invalid object type")
	ErrObjectTypeNotFound     = cerr.NewAsertoError("E21002", codes.NotFound, http.StatusNotFound, "object type not found")
	ErrInvalidRelationType    = cerr.NewAsertoError("E21003", codes.InvalidArgument, http.StatusBadRequest, "invalid relation type")
	ErrRelationTypeNotFound   = cerr.NewAsertoError("E21004", codes.NotFound, http.StatusNotFound, "relation type not found")
	ErrInvalidPermissionType  = cerr.NewAsertoError("E21005", codes.InvalidArgument, http.StatusBadRequest, "invalid permission type")
	ErrPermissionTypeNotFound = cerr.NewAsertoError("E21006", codes.NotFound, http.StatusNotFound, "permission type not found")
	ErrOperationNotFound      = cerr.NewAsertoError("E21007", codes.NotFound, http.StatusNotFound, "operation not found")
)
