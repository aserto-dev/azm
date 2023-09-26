package cache

import (
	"github.com/aserto-dev/azm/paging"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
)

type ObjectTypeSlice []*dsc2.ObjectType
type RelationTypeSlice []*dsc2.RelationType
type PermissionSlice []*dsc2.Permission

func (s ObjectTypeSlice) Paginate(page *dsc2.PaginationRequest) (*paging.Result[*dsc2.ObjectType], error) {
	if paging.IsCountOnly(page) {
		return &paging.Result[*dsc2.ObjectType]{
			PageInfo: &dsc2.PaginationResponse{ResultSize: int32(len(s))},
		}, nil
	}

	return paging.PaginateSlice(
		s,
		page,
		1,
		func(keys []string, ot *dsc2.ObjectType) bool { return keys[0] == ot.Name },
		func(ot *dsc2.ObjectType) []string { return []string{ot.Name} },
	)
}

func (s RelationTypeSlice) Paginate(page *dsc2.PaginationRequest) (*paging.Result[*dsc2.RelationType], error) {
	if paging.IsCountOnly(page) {
		return &paging.Result[*dsc2.RelationType]{
			PageInfo: &dsc2.PaginationResponse{ResultSize: int32(len(s))},
		}, nil
	}

	return paging.PaginateSlice(
		s,
		page,
		2,
		func(keys []string, relType *dsc2.RelationType) bool {
			return keys[0] == relType.ObjectType && keys[1] == relType.Name
		},
		func(relType *dsc2.RelationType) []string { return []string{relType.ObjectType, relType.Name} },
	)
}

func (s PermissionSlice) Paginate(page *dsc2.PaginationRequest) (*paging.Result[*dsc2.Permission], error) {
	if paging.IsCountOnly(page) {
		return &paging.Result[*dsc2.Permission]{
			PageInfo: &dsc2.PaginationResponse{ResultSize: int32(len(s))},
		}, nil
	}

	return paging.PaginateSlice(
		s,
		page,
		1,
		func(keys []string, p *dsc2.Permission) bool { return keys[0] == p.Name },
		func(p *dsc2.Permission) []string { return []string{p.Name} },
	)
}
