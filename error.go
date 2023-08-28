package azm

import (
	"net/http"

	cerr "github.com/aserto-dev/errors"

	"google.golang.org/grpc/codes"
)

var (
	ErrModelNodeNotFound          = cerr.NewAsertoError("E10000", codes.NotFound, http.StatusNotFound, "manifest does not contain a model node")
	ErrSchemaVersionNotFound      = cerr.NewAsertoError("E10001", codes.NotFound, http.StatusNotFound, "manifest does not contain a model.version field")
	ErrInvalidSchemaVersion       = cerr.NewAsertoError("E10002", codes.InvalidArgument, http.StatusBadRequest, "invalid or unsupported schema version")
	ErrInvalidRelationDefinitions = cerr.NewAsertoError("E10003", codes.InvalidArgument, http.StatusBadRequest, "model expected array of relation definitions")
	ErrInvalidRelationDefinition  = cerr.NewAsertoError("E10003", codes.InvalidArgument, http.StatusBadRequest, "model expected map of relation definition")
)
