// nolint: staticcheck
package migrate_test

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/aserto-dev/azm/migrate"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/samber/lo"
)

var (
	// base directory object types.
	ObjectTypes = []*dsc2.ObjectType{
		{Name: "system", DisplayName: "System", IsSubject: false, Ordinal: 900, Status: uint32(dsc2.Flag_FLAG_HIDDEN | dsc2.Flag_FLAG_SYSTEM)},
		{Name: "user", DisplayName: "User", IsSubject: true, Ordinal: 100, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{Name: "identity", DisplayName: "Identity", IsSubject: false, Ordinal: 300, Status: uint32(dsc2.Flag_FLAG_SYSTEM | dsc2.Flag_FLAG_READONLY)},
		{Name: "group", DisplayName: "Group", IsSubject: true, Ordinal: 200, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{Name: "application", DisplayName: "Application", IsSubject: false, Ordinal: 400, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{Name: "resource", DisplayName: "Resource", IsSubject: false, Ordinal: 500, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{Name: "user-v1", DisplayName: "UserV1", IsSubject: true, Ordinal: 1000, Status: uint32(dsc2.Flag_FLAG_HIDDEN | dsc2.Flag_FLAG_SYSTEM | dsc2.Flag_FLAG_SHADOW | dsc2.Flag_FLAG_READONLY)},
	}

	// base directory relation types.
	RelationTypes = []*dsc2.RelationType{
		{ObjectType: "system", Name: "user", DisplayName: "system:user", Ordinal: 900, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{ObjectType: "identity", Name: "identifier", DisplayName: "identity:identifier", Ordinal: 200, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{ObjectType: "group", Name: "member", DisplayName: "group:member", Ordinal: 100, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{ObjectType: "application", Name: "user", DisplayName: "application:user", Ordinal: 400, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
		{ObjectType: "user", Name: "manager", DisplayName: "user:manager", Ordinal: 300, Status: uint32(dsc2.Flag_FLAG_SYSTEM)},
	}
)

func TestMigrateBaseToAnnotatedModelV3(t *testing.T) {
	sort.Slice(ObjectTypes, func(i, j int) bool {
		return ObjectTypes[i].Ordinal < ObjectTypes[j].Ordinal
	})

	md := &migrate.Metadata{
		ObjectTypes:   ObjectTypes,
		RelationTypes: RelationTypes,
		Permissions:   []*dsc2.Permission{},
	}

	w := os.Stderr
	WriteManifestHeader(w)
	WriteManifestInfo(w, "manifest-v2-gen.yaml", "automatic migration of v2 model to annotated v3 manifest")
	WriteManifestModel(w, 2, 3)
	WriteManifestTypes(w)

	for _, ot := range md.ObjectTypes {
		WriteTypeInstance(w, 2, ot)

		rts := lo.Filter(RelationTypes, func(item *dsc2.RelationType, index int) bool {
			return item.ObjectType == ot.Name
		})

		if len(rts) > 0 {
			WriteRelations(w, 4)
		}

		for _, rt := range rts {
			WriteRelationInstance(w, 6, rt)
		}

		fmt.Fprintln(w)
	}
}

func WriteManifestHeader(w io.Writer) {
	fmt.Fprint(w, "# yaml-language-server: $schema=https://www.topaz.sh/schema/manifest.json\n")
	fmt.Fprint(w, "---\n")
	fmt.Fprintln(w)
}

func WriteManifestInfo(w io.Writer, filename, description string) {
	fmt.Fprintf(w, "### filename: %s ###\n", filename)
	fmt.Fprintf(w, "### description: %s ###\n", description)
	fmt.Fprintln(w)
}

func WriteManifestModel(w io.Writer, indent, version int) {
	fmt.Fprint(w, "### model ###\n")
	fmt.Fprint(w, "model:\n")
	fmt.Fprintf(w, "%sversion: %d\n", space(indent), version)
	fmt.Fprintln(w)
}

func WriteManifestTypes(w io.Writer) {
	fmt.Fprint(w, "### object type definitions ###\n")
	fmt.Fprint(w, "types:\n")
}

func WriteTypeInstance(w io.Writer, indent int, instance *dsc2.ObjectType) {
	fmt.Fprintf(w, "%s### %s: %s ###\n", space(indent), "display_name", instance.DisplayName)
	fmt.Fprintf(w, "%s### %s: %d ###\n", space(indent), "ordinal", instance.Ordinal)
	fmt.Fprintf(w, "%s%s:\n", space(indent), instance.Name)
}

func WriteRelations(w io.Writer, indent int) {
	fmt.Fprintf(w, "%srelations:\n", space(indent))
}

func WriteRelationInstance(w io.Writer, indent int, instance *dsc2.RelationType) {
	target := "user"

	fmt.Fprintf(w, "%s### %s: %s ###\n", space(indent), "display_name", instance.DisplayName)
	fmt.Fprintf(w, "%s### %s: %d ###\n", space(indent), "ordinal", instance.Ordinal)
	fmt.Fprintf(w, "%s%s: %s\n", space(indent), instance.Name, target)
}

func space(count int) string {
	return strings.Repeat(" ", count)
}
