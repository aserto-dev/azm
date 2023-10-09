package migrate

import dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"

type Metadata struct {
	ObjectTypes   []*dsc2.ObjectType
	RelationTypes []*dsc2.RelationType
	Permissions   []*dsc2.Permission
}
