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
	"github.com/samber/lo"

	schema2 "entgo.io/ent/schema"
)

type importer struct {
	DSN               string
	SchemaDir         string
	SkipFields        []string
	AddImports        []string
	FieldProvider     func(table string, column string) any
	AstRewriteHandler func(src []byte) ([]byte, error)
}

func NewImporter(
	dsn string,
	dir string,
	skipFields []string,
	addImports []string,
	provider func(string, string) any,
	rewriter func(src []byte) ([]byte, error),
) *importer {
	return &importer{
		DSN:               dsn,
		SchemaDir:         dir,
		SkipFields:        skipFields,
		AddImports:        addImports,
		FieldProvider:     provider,
		AstRewriteHandler: rewriter,
	}
}

func (i *importer) FieldBuilder(column *schema.Column) any {
	switch column.Type.Type.(type) {
	case *schema.StringType:
		return field.String
	case *schema.IntegerType:
		return field.Int64
	case *schema.FloatType:
		return field.Float
	case *schema.DecimalType:
		return field.Float
	case *schema.BoolType:
		return field.Bool
	case *schema.TimeType:
		return field.String
	case *schema.BinaryType:
		return field.Bytes
	default:
		panic("not impl")
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
			Annotations: []schema2.Annotation{
				entsql.Annotation{
					Table: table.Name,
				},
			},
		}

		for _, column := range table.Columns {
			var (
				column  = column
				builder *fieldBuilder
			)

			if i.SkipFields != nil && lo.Contains(i.SkipFields, column.Name) {
				continue
			}

			if i.FieldProvider == nil || i.FieldProvider(table.Name, column.Name) == nil {
				builder = NewFieldBuilder(column.Name, i.FieldBuilder(column))
			} else {
				builder = NewFieldBuilder(column.Name, i.FieldProvider(table.Name, column.Name))
			}

			mutator.Fields = append(
				mutator.Fields,

				builder.Default(column.Default).
					SchemaType(map[string]string{
						dialect.MySQL: column.Type.Raw, // Override MySQL.
					}).
					StructTag(newStructTag(column.Name)).
					Comment(column.Attrs),
			)
		}

		mutations = append(mutations, mutator)
	}

	err = schemast.Mutate(ctx, mutations...)
	if err := ctx.Print(i.SchemaDir); err != nil {
		panic(fmt.Errorf("failed: %v", err))
	}

	dstfmt.FmtDir(i.SchemaDir, i.AstRewriteHandler, i.AddImports)
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
