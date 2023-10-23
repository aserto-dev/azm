package migrate

type ObjRelSub struct {
	Object   string
	Relation string
	Subject  string
}

type ObjRelSubContainer struct {
	count int
	index map[string]map[string]map[string]struct{}
}

func NewObjRelSubContainer() *ObjRelSubContainer {
	return &ObjRelSubContainer{
		index: make(map[string]map[string]map[string]struct{}),
	}
}

func (o *ObjRelSubContainer) Add(ors *ObjRelSub) {
	obj, ok := o.index[ors.Object]
	if !ok {
		o.index[ors.Object] = make(map[string]map[string]struct{})
		obj = o.index[ors.Object]
	}

	rel, ok := obj[ors.Relation]
	if !ok {
		obj[ors.Relation] = make(map[string]struct{})
		rel = obj[ors.Relation]
	}

	if _, ok := rel[ors.Subject]; !ok {
		rel[ors.Subject] = struct{}{}
	}
	o.count++
}

func (o *ObjRelSubContainer) Get(obj, rel, sub string) (*ObjRelSub, bool) {
	rels, ok := o.index[obj]
	if !ok {
		return nil, false
	}
	subs, ok := rels[rel]
	if !ok {
		return nil, false
	}
	if _, ok := subs[sub]; !ok {
		return nil, false
	}
	return &ObjRelSub{Object: obj, Relation: rel, Subject: sub}, true
}

func (o *ObjRelSubContainer) GetRels(obj string) []*ObjRelSub {
	results := []*ObjRelSub{}
	rels, ok := o.index[obj]
	if !ok {
		return results
	}
	for rn := range rels {
		results = append(results, &ObjRelSub{Object: obj, Relation: rn, Subject: ""})
	}
	return results
}

func (o *ObjRelSubContainer) GetSubs(obj, rel string) []*ObjRelSub {
	results := []*ObjRelSub{}
	rels, ok := o.index[obj]
	if !ok {
		return results
	}
	subs, ok := rels[rel]
	if !ok {
		return results
	}
	for sn := range subs {
		results = append(results, &ObjRelSub{Object: obj, Relation: rel, Subject: sn})
	}
	return results
}

func (o *ObjRelSubContainer) All() []*ObjRelSub {
	results := []*ObjRelSub{}
	for on, rels := range o.index {
		for rn, subs := range rels {
			for sn := range subs {
				results = append(results, &ObjRelSub{
					Object:   on,
					Relation: rn,
					Subject:  sn,
				})
			}
		}
	}
	return results
}

func (o *ObjRelSubContainer) RelationCount() int {
	return o.count
}
