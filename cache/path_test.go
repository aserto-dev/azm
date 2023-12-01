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

func (pm PathMap) GetPath(or *model.ObjectRelation) []*model.ObjectRelation {
	if or == nil {
		return []*model.ObjectRelation{}
	}

	p1, ok := pm[or.Object]
	if !ok {
		return []*model.ObjectRelation{}
	}

	p2, ok := p1[or.Relation]
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

	pm := createPathMap(m)
	require.NotNil(t, pm)

	// plot all paths for all roots.
	roots := []*model.ObjectRelation{}
	for on, rns := range *pm {
		for rn := range rns {
			roots = append(roots, model.NewObjectRelation(on, rn))
		}
	}

	for i := 0; i < len(roots); i++ {
		path := pm.WalkPath(roots[i], []string{})
		fmt.Println(strings.Join(path, " -> "))
	}
}

func (pm PathMap) WalkPath(or *model.ObjectRelation, path []string) []string {
	paths := pm.GetPath(or)
	for i := 0; i < len(paths); i++ {
		path = append(path, paths[i].String())
		pm.WalkPath(paths[i], path)
	}
	return path
}

func (pm PathMap) plotPaths(w io.Writer, or *model.ObjectRelation) {
	paths := pm.GetPath(or)

	for _, p := range paths {
		fmt.Fprintf(w, "%s:%s ", or.Object, or.Relation)

		fmt.Fprintf(w, "-> %s:%s ", p.Object, p.Relation)

		for _, v := range pm.GetPath(p) {
			fmt.Fprintf(w, "-> %s:%s ", v.Object, v.Relation)
		}

		fmt.Fprintf(w, "\n")
	}
}

func createPathMap(m *model.Model) *PathMap {
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

	return &pm
}

func expandPerm(m *model.Model, on model.ObjectName, pn model.PermissionName) []*model.ObjectRelation {
	result := []*model.ObjectRelation{}

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

func expandRel(m *model.Model, on model.ObjectName, rn model.RelationName) []*model.ObjectRelation {
	result := []*model.ObjectRelation{}

	relations, ok := m.Objects[on].Relations[rn]
	if !ok {
		return result
	}

	for _, r := range relations {
		if r.Direct != "" {
			result = append(result, &model.ObjectRelation{
				Object:   r.Direct,
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
				Object:   r.Wildcard,
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
						Object:   rel.Direct,
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
						Object:   rel.Wildcard,
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