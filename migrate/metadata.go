package migrate

import dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"

type Metadata struct {
	ObjectTypes   []*dsc2.ObjectType
	RelationTypes []*dsc2.RelationType
	Permissions   []*dsc2.Permission
}

const Obsolete int32 = 16

var RefObjectTypes = map[string]*dsc2.ObjectType{
	"system":      {Name: "system", DisplayName: "System", IsSubject: false, Ordinal: 6, Status: uint32(dsc2.Flag_FLAG_HIDDEN | dsc2.Flag_FLAG_SYSTEM | dsc2.Flag(Obsolete))},
	"user":        {Name: "user", DisplayName: "User", IsSubject: true, Ordinal: 1, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
	"identity":    {Name: "identity", DisplayName: "Identity", IsSubject: false, Ordinal: 2, Status: uint32(dsc2.Flag_FLAG_SYSTEM | dsc2.Flag_FLAG_READONLY)},
	"group":       {Name: "group", DisplayName: "Group", IsSubject: true, Ordinal: 3, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
	"application": {Name: "application", DisplayName: "Application", IsSubject: false, Ordinal: 4, Status: uint32(dsc2.Flag_FLAG_SYSTEM | dsc2.Flag(Obsolete))},
	"resource":    {Name: "resource", DisplayName: "Resource", IsSubject: false, Ordinal: 5, Status: uint32(dsc2.Flag_FLAG_SYSTEM | dsc2.Flag(Obsolete))},
	"user-v1":     {Name: "user-v1", DisplayName: "UserV1", IsSubject: true, Ordinal: 0, Status: uint32(dsc2.Flag_FLAG_HIDDEN | dsc2.Flag_FLAG_SYSTEM | dsc2.Flag_FLAG_SHADOW | dsc2.Flag_FLAG_READONLY | dsc2.Flag(Obsolete))},
}

var RefRelationTypes = map[string]map[string]*dsc2.RelationType{
	"system":      {"user": {ObjectType: "system", Name: "user", DisplayName: "system#user", Ordinal: 0, Status: uint32(dsc2.Flag_FLAG_SYSTEM | dsc2.Flag(Obsolete))}},
	"identity":    {"identifier": {ObjectType: "identity", Name: "identifier", Unions: []string{"user"}, DisplayName: "identity#identifier", Ordinal: 0, Status: uint32(dsc2.Flag_FLAG_SYSTEM)}},
	"group":       {"member": {ObjectType: "group", Name: "member", Unions: []string{"user"}, DisplayName: "group#member", Ordinal: 0, Status: uint32(dsc2.Flag_FLAG_SYSTEM)}},
	"application": {"user": {ObjectType: "application", Name: "user", DisplayName: "application#user", Ordinal: 0, Status: uint32(dsc2.Flag_FLAG_SYSTEM | dsc2.Flag(Obsolete))}},
	"user":        {"manager": {ObjectType: "user", Name: "manager", Unions: []string{"user"}, DisplayName: "user#manager", Ordinal: 0, Status: uint32(dsc2.Flag_FLAG_SYSTEM)}},
}

func IsObsolete(status dsc2.Flag) bool {
	return status&dsc2.Flag(Obsolete) != 0
}
