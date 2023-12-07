package cache_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/stretchr/testify/require"
)

type PathMap map[model.ObjectName]map[model.RelationName][]*model.RelationRef

func (pm PathMap) GetPath(or *model.RelationRef) []*model.RelationRef {
	if or == nil {
		return []*model.RelationRef{}
	}

	p1, ok := pm[or.Object]
	if !ok {
		return []*model.RelationRef{}
	}

	p2, ok := p1[or.Relation]
	if !ok {
		return []*model.RelationRef{}
	}

	return p2
}

// func walkPath(m *model.Model, rr *model.RelationRef, path []string) []string {

// }

func TestPathMap(t *testing.T) {
	r, err := os.Open("./path_test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)

	m, err := v3.Load(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	c := cache.New(m)
	require.NotNil(t, c)

	pm := createPathMap(m)
	require.NotNil(t, pm)

	// plot all paths for all roots.
	for on, rns := range *pm {
		for rn := range rns {
			path := pm.WalkPath(model.NewRelationRef(on, rn), []string{})
			fmt.Printf("%s:%s: %s\n", on, rn, strings.Join(path, " -> "))
		}
	}
}

func (pm PathMap) WalkPath(or *model.RelationRef, path []string) []string {
	paths := pm.GetPath(or)
	for i := 0; i < len(paths); i++ {
		path = append(path, paths[i].String())
		pm.WalkPath(paths[i], path)
	}
	return path
}

func createPathMap(m *model.Model) *PathMap {
	pm := PathMap{}

	// create roots
	for on, o := range m.Objects {
		if _, ok := pm[on]; !ok {
			pm[on] = map[model.RelationName][]*model.RelationRef{}
		}

		p1 := pm[on]

		for pn := range o.Permissions {
			if _, ok := p1[pn.RN()]; !ok {
				p1[pn.RN()] = []*model.RelationRef{}
			}
		}

		for rn := range o.Relations {
			if _, ok := p1[rn]; !ok {
				p1[rn] = []*model.RelationRef{}
			}
		}
	}

	// set leaves
	for on, o := range m.Objects {
		p1 := pm[on]

		for pn := range o.Permissions {
			p1[pn.RN()] = expandPerm(m, on, pn)
		}

		for rn := range o.Relations {
			p1[rn] = expandRel(m, on, rn)
		}
	}

	return &pm
}

func expandPerm(m *model.Model, on model.ObjectName, pn model.PermissionName) []*model.RelationRef {
	result := []*model.RelationRef{}

	p, ok := m.Objects[on].Permissions[pn]
	if !ok {
		return result
	}

	for _, r := range p.Union {
		result = append(result, resolve(m, on, model.RelationName(r.RelOrPerm)))
	}

	for range p.Intersection {
		panic("not implemented")
	}

	if p.Exclusion != nil {
		panic("not implemented")
	}

	return result
}

func expandRel(m *model.Model, on model.ObjectName, rn model.RelationName) []*model.RelationRef {
	result := []*model.RelationRef{}

	relation, ok := m.Objects[on].Relations[rn]
	if !ok {
		return result
	}

	for _, r := range relation.Union {
		if r.Direct != nil {
			result = append(result, r.Direct)
		}

		if r.Subject != nil {
			result = append(result, &model.RelationRef{
				Object:   r.Subject.Object,
				Relation: r.Subject.Relation,
			})
		}

		if r.Wildcard != nil {
			result = append(result, r.Wildcard)
		}
	}

	return result
}

func resolve(m *model.Model, on model.ObjectName, rn model.RelationName) *model.RelationRef {
	if strings.Contains(rn.String(), v3.ArrowIdentifier) {
		parts := strings.Split(rn.String(), v3.ArrowIdentifier)

		rn = model.RelationName(parts[0])

		if _, ok := m.Objects[on].Relations[rn]; ok { // 	if c.RelationExists(on, rn) {
			for _, rel := range m.Objects[on].Relations[rn].Union {
				if rel.Direct != nil {
					return rel.Direct
				}

				if rel.Subject != nil {
					return rel.Subject.RelationRef
				}

				if rel.Wildcard != nil {
					return rel.Wildcard
				}
			}
		}
	}

	return &model.RelationRef{
		Object:   on,
		Relation: rn,
	}
}
