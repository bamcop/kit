package kent_import

import (
	"context"
	"fmt"
	"strings"

	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlclient"
	"entgo.io/contrib/schemast"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"github.com/bamcop/kit/ent_util/dstfmt"
	"github.com/iancoleman/strcase"

	ents "entgo.io/ent/schema"
)

type importer struct {
	DSN           string
	SchemaDir     string
	FieldProvider func(table string, column string) any
}

func NewImporter(dsn string, dir string, provider func(string, string) any) *importer {
	return &importer{
		DSN:           dsn,
		SchemaDir:     dir,
		FieldProvider: provider,
	}
}

func (i *importer) Execute() error {
	info := inspectInfo(i.DSN)

	ctx, err := schemast.Load(i.SchemaDir)
	if err != nil {
		panic(fmt.Errorf("failed: %v", err))
	}

	var mutations []schemast.Mutator
	for _, table := range info.Tables {
		table := table

		mutator := &schemast.UpsertSchema{
			Name: strcase.ToCamel(strings.TrimLeft(table.Name, "t_")),
			Annotations: []ents.Annotation{
				entsql.Annotation{
					Table: table.Name,
				},
				entsql.WithComments(true),
			},
		}

		for _, column := range table.Columns {
			var (
				column  = column
				builder *fieldBuilder
			)

			switch column.Type.Raw {
			case "string":
				builder = NewFieldBuilder(column.Name, field.String)

			}

			mutator.Fields = append(
				mutator.Fields,

				builder.SchemaType(map[string]string{
					dialect.MySQL: column.Type.Raw, // Override MySQL.
				}).
					Default(column.Default).
					StructTag(newStructTag(column.Name)),
				//Comment(column.Attrs)

				//newFieldBuilder(i.FieldProvider(table.Name, column.Name)).
				//	SchemaType(map[string]string{
				//		dialect.MySQL: column.Type.Raw, // Override MySQL.
				//	}).
				//	Default(column.Default).
				//	StructTag(newStructTag(column.Name)).
				//	Comment(column.Attrs),
			)
		}

		mutations = append(mutations, mutator)
	}

	err = schemast.Mutate(ctx, mutations...)
	if err := ctx.Print(i.SchemaDir); err != nil {
		panic(fmt.Errorf("failed: %v", err))
	}

	dstfmt.FmtDir(i.SchemaDir)
	return nil
}

func inspectInfo(dsn string) *schema.Schema {
	client, err := sqlclient.Open(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	info, err := client.InspectSchema(context.Background(), "", nil)
	if err != nil {
		panic(err)
	}

	return info
}

func newStructTag(str string) string {
	return fmt.Sprintf("json:\"%s\"", str)
}
