package migrate

type ObjPermRel struct {
	Object     string
	Permission string
	Relation   string
}

type ObjPermRelContainer struct {
	index map[string]map[string]map[string]struct{}
}

func NewObjPermRelContainer() *ObjPermRelContainer {
	return &ObjPermRelContainer{
		index: make(map[string]map[string]map[string]struct{}),
	}
}

func (o *ObjPermRelContainer) Add(opr *ObjPermRel) {
	obj, ok := o.index[opr.Object]
	if !ok {
		o.index[opr.Object] = make(map[string]map[string]struct{})
		obj = o.index[opr.Object]
	}

	perm, ok := obj[opr.Permission]
	if !ok {
		obj[opr.Permission] = make(map[string]struct{})
		perm = obj[opr.Permission]
	}

	if _, ok := perm[opr.Relation]; !ok {
		perm[opr.Relation] = struct{}{}
	}
}

func (o *ObjPermRelContainer) GetPerms(obj string) map[string]map[string]struct{} {
	perms, ok := o.index[obj]
	if !ok {
		return map[string]map[string]struct{}{}
	}
	return perms
}

func (o *ObjPermRelContainer) All() []*ObjPermRel {
	results := []*ObjPermRel{}
	for on, perms := range o.index {
		for pn, rels := range perms {
			for rn := range rels {
				results = append(results, &ObjPermRel{
					Object:     on,
					Permission: pn,
					Relation:   rn,
				})
			}
		}
	}
	return results
}
