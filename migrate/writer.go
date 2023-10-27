package migrate

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aserto-dev/azm/model"
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/samber/lo"
)

func WriteManifest(w io.Writer, md *Metadata, pts *ObjPermRelContainer) {

	writeManifestHeader(w)
	writeManifestInfo(w, "manifest.yaml", "automatic migration of v2 model to annotated v3 manifest")
	writeManifestModel(w, 2, 3)
	writeManifestTypes(w)

	for _, ot := range md.ObjectTypes {
		rts := lo.Filter(md.RelationTypes, func(item *dsc2.RelationType, index int) bool {
			return item.ObjectType == ot.Name
		})

		if len(rts) == 0 && len(pts.GetPerms(ot.Name)) == 0 {
			writeTypeInstance(w, 2, ot, true)
			fmt.Fprintln(w)
			continue
		}

		writeTypeInstance(w, 2, ot, false)

		if len(rts) > 0 {
			writeRelations(w, 4)
		}

		for _, rt := range rts {
			writeRelationInstance(w, 6, rt)
		}

		if len(pts.GetPerms(ot.Name)) > 0 {
			writePermissions(w, 4)
		}

		writePermissionInstance(w, 6, pts.GetPerms(ot.Name))

		fmt.Fprintln(w)
	}
}

func writeManifestHeader(w io.Writer) {
	fmt.Fprint(w, "# yaml-language-server: $schema=https://www.topaz.sh/schema/manifest.json\n")
	fmt.Fprint(w, "---\n")
	fmt.Fprintln(w)
}

func writeManifestInfo(w io.Writer, filename, description string) {
	fmt.Fprintf(w, "### filename: %s ###\n", filename)
	fmt.Fprintf(w, "### datetime: %s ###\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "### description: %s ###\n", description)
	fmt.Fprintln(w)
}

func writeManifestModel(w io.Writer, indent, version int) {
	fmt.Fprint(w, "### model ###\n")
	fmt.Fprint(w, "model:\n")
	fmt.Fprintf(w, "%sversion: %d\n", space(indent), version)
	fmt.Fprintln(w)
}

func writeManifestTypes(w io.Writer) {
	fmt.Fprint(w, "### object type definitions ###\n")
	fmt.Fprint(w, "types:\n")
}

func writeTypeInstance(w io.Writer, indent int, instance *dsc2.ObjectType, this bool) {
	fmt.Fprintf(w, "%s### %s: %s ###\n", space(indent), "display_name", instance.DisplayName)

	name, err := model.NormalizeIdentifier(instance.Name)

	if err != nil {
		fmt.Fprintf(w, "%s%s:\n", space(indent), "### INVALID IDENTIFIER >>>")
		fmt.Fprintf(w, "%s%s:\n", space(indent), "### "+name)
		fmt.Fprintf(w, "%s%s:\n", space(indent), "### INVALID IDENTIFIER <<<")
		return
	}

	fmt.Fprintf(w, "%s%s:%s\n", space(indent), name, iff(this, " {}", ""))
}

func writeRelations(w io.Writer, indent int) {
	fmt.Fprintf(w, "%srelations:\n", space(indent))
}

func writeRelationInstance(w io.Writer, indent int, instance *dsc2.RelationType) {
	fmt.Fprintf(w, "%s### %s: %s ###\n", space(indent), "display_name", instance.DisplayName)

	name, err := model.NormalizeIdentifier(instance.Name)
	if err != nil {
		fmt.Fprintf(w, "%s%s:\n", space(indent), "### INVALID IDENTIFIER >>>")
		fmt.Fprintf(w, "%s%s: %s\n", space(indent), "### "+name, strings.Join(instance.Unions, " | "))
		fmt.Fprintf(w, "%s%s:\n", space(indent), "### INVALID IDENTIFIER <<<")
	}

	fmt.Fprintf(w, "%s%s: %s\n", space(indent), name, strings.Join(instance.Unions, " | "))
}

func writePermissions(w io.Writer, indent int) {
	fmt.Fprintf(w, "%spermissions:\n", space(indent))
}

func writePermissionInstance(w io.Writer, indent int, instances map[string]map[string]struct{}) {
	for pn, rels := range instances {
		relArray := lo.MapToSlice(rels, func(k string, v struct{}) string {
			return k
		})
		fmt.Fprintf(w, "%s### %s: %s ###\n", space(indent), "display_name", pn)

		name, err := model.NormalizeIdentifier(pn)
		if err != nil {
			fmt.Fprintf(w, "%s%s:\n", space(indent), "### INVALID IDENTIFIER >>>")
			fmt.Fprintf(w, "%s%s: %s\n", space(indent), "### "+pn, strings.Join(relArray, " | "))
			fmt.Fprintf(w, "%s%s:\n", space(indent), "### INVALID IDENTIFIER <<<")
		}

		fmt.Fprintf(w, "%s%s: %s\n", space(indent), name, strings.Join(relArray, " | "))
	}
}

func space(count int) string {
	return strings.Repeat(" ", count)
}

func iff[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}
