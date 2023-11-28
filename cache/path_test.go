package cache_test

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aserto-dev/azm/cache"
	"github.com/aserto-dev/azm/model"
	v3 "github.com/aserto-dev/azm/v3"
	"github.com/stretchr/testify/require"
)

type PathMap map[model.ObjectName]map[model.RelationName][]*model.ObjectRelation

func (pm PathMap) GetPath(on model.ObjectName, rn model.RelationName) []*model.ObjectRelation {
	p1, ok := pm[on]
	if !ok {
		return []*model.ObjectRelation{}
	}

	p2, ok := p1[rn]
	if !ok {
		return []*model.ObjectRelation{}
	}

	return p2
}

func TestPathMap(t *testing.T) {
	r, err := os.Open("./path_test.yaml")
	require.NoError(t, err)
	require.NotNil(t, r)

	m, err := v3.Load(r)
	require.NoError(t, err)
	require.NotNil(t, m)

	c := cache.New(m)
	require.NotNil(t, c)

	pm, err := createPathMap(m)
	require.NoError(t, err)
	require.NotNil(t, pm)

	// plot all paths for all roots.
	for on, rns := range *pm {
		for rn := range rns {
			pm.plotPaths(os.Stderr, on, rn)
		}
	}
}

func (pm PathMap) plotPaths(w io.Writer, on model.ObjectName, rn model.RelationName) {
	paths := pm.GetPath(on, rn)

	for _, p := range paths {
		fmt.Fprintf(w, "%s:%s ", on, rn)

		fmt.Fprintf(w, "-> %s:%s ", p.Object, p.Relation)

		for _, v := range pm.GetPath(p.Object, p.Relation) {
			fmt.Fprintf(w, "-> %s:%s ", v.Object, v.Relation)
		}

		fmt.Fprintf(w, "\n")
	}
}

func createPathMap(m *model.Model) (*PathMap, error) {
	pm := PathMap{}

	// create roots
	for on, o := range m.Objects {
		if _, ok := pm[on]; !ok {
			pm[on] = map[model.RelationName][]*model.ObjectRelation{}
		}

		p1 := pm[on]

		for pn := range o.Permissions {
			if _, ok := p1[pn.RN()]; !ok {
				p1[pn.RN()] = []*model.ObjectRelation{}
			}
		}

		for rn := range o.Relations {
			if _, ok := p1[rn]; !ok {
				p1[rn] = []*model.ObjectRelation{}
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

	return &pm, nil
}

func expandPerm(m *model.Model, on model.ObjectName, pn model.PermissionName) []*model.ObjectRelation {
	result := []*model.ObjectRelation{}

	p, ok := m.Objects[on].Permissions[pn]
	if !ok {
		return result
	}

	for _, r := range p.Union {
		result = append(result, resolve(m, on, model.RelationName(r)))
	}

	for _, _ = range p.Intersection {
		panic("not implemented")
	}

	if p.Exclusion != nil {
		panic("not implemented")
	}

	if p.Arrow != nil {
		panic("not implemented")
	}

	return result
}

func expandRel(m *model.Model, on model.ObjectName, rn model.RelationName) []*model.ObjectRelation {
	result := []*model.ObjectRelation{}

	relations, ok := m.Objects[on].Relations[rn]
	if !ok {
		return result
	}

	for _, r := range relations {
		if r.Direct != "" {
			result = append(result, &model.ObjectRelation{
				Object:   model.ObjectName(r.Direct),
				Relation: "",
			})
		}

		if r.Subject != nil {
			result = append(result, &model.ObjectRelation{
				Object:   r.Subject.Object,
				Relation: r.Subject.Relation,
			})
		}

		if r.Wildcard != "" {
			result = append(result, &model.ObjectRelation{
				Object:   model.ObjectName(r.Wildcard),
				Relation: "*",
			})
		}
	}

	return result
}

func resolve(m *model.Model, on model.ObjectName, rn model.RelationName) *model.ObjectRelation {
	if strings.Contains(rn.String(), v3.ArrowIdentifier) {
		parts := strings.Split(rn.String(), v3.ArrowIdentifier)

		rn := model.RelationName(parts[0])

		if _, ok := m.Objects[on].Relations[rn]; ok { // 	if c.RelationExists(on, rn) {
			for _, rel := range m.Objects[on].Relations[rn] {
				if rel.Direct != "" {
					return &model.ObjectRelation{
						Object:   model.ObjectName(rel.Direct),
						Relation: model.RelationName(parts[1]),
					}
				}

				if rel.Subject != nil {
					return &model.ObjectRelation{
						Object:   rel.Subject.Object,
						Relation: rel.Subject.Relation,
					}
				}

				if rel.Wildcard != "" {
					return &model.ObjectRelation{
						Object:   model.ObjectName(rel.Wildcard),
						Relation: "*",
					}
				}
			}
		}
	}

	return &model.ObjectRelation{
		Object:   on,
		Relation: rn,
	}
}
