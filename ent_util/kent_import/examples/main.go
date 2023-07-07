package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bamcop/kit/ast_util/aq"
	"github.com/bamcop/kit/ent_util/kent_import"
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
	"github.com/helloyi/goastch/goastcher"
	"github.com/samber/lo"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	TmplMinxin = `
// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		ent_util.BaseMixin{},
	}
}
`
	TmplAnnotation = `
// Mixin of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "t_user"},
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
`
)

var (
	WorkingDir string
)

func init() {
	WorkingDir, _ = os.Getwd()
}

func main() {
	rewriter := func(src []byte) ([]byte, error) {
		doc, err := aq.New("", src)
		if err != nil {
			panic(err)
		}

		var (
			recv  string
			table string
		)

		doc.Find(
			goastcher.Has(goastcher.FuncDecl(goastcher.Anything())).Bind("1"),
			goastcher.Has(goastcher.HasName(goastcher.Equals("Annotations"))).Bind("2"),
		).ForEach(func(root dst.Node, node dst.Node) {
			recv = node.(*dst.FuncDecl).Recv.List[0].Type.(*dst.Ident).Name
		}).FindLeave(
			goastcher.HasDescendant(goastcher.MatchCode(`entsql.Annotation{Table:`)).Bind("3"),
		).Find(
			goastcher.HasDescendant(goastcher.KeyValueExpr(goastcher.Anything())).Bind("8"),
		).ForEach(func(root dst.Node, node dst.Node) {
			table = node.(*dst.KeyValueExpr).Value.(*dst.BasicLit).Value
			table = strings.Trim(table, `"`)
		}).DebugPrint("最终结果")

		dstutil.Apply(doc.Root(), func(cursor *dstutil.Cursor) bool {
			if cursor.Node() == doc.Root() {
				str1 := strings.ReplaceAll(TmplMinxin, "User", recv)
				str2 := strings.ReplaceAll(TmplAnnotation, "User", recv)
				str2 = strings.ReplaceAll(str2, "t_user", table)

				f1 := aq.Gen.MustNewFunc(str1)
				f2 := aq.Gen.MustNewFunc(str2)

				doc.Root().Decls = lo.Filter(doc.Root().Decls, func(item dst.Decl, index int) bool {
					v, ok := item.(*dst.FuncDecl)
					if !ok {
						return true
					}
					if v.Name.Name != "Annotations" {
						return true
					}
					return false
				})

				doc.Root().Decls = append(doc.Root().Decls, f1.(dst.Decl))
				doc.Root().Decls = append(doc.Root().Decls, f2.(dst.Decl))

				return false
			}
			return true
		}, nil)

		doc.DebugPrintRoot()

		return []byte(doc.RootString()), nil
	}

	var (
		// DSN的格式: https://atlasgo.io/concepts/url
		dsn = fmt.Sprintf("sqlite://%s/sample.db", WorkingDir)
		dir = filepath.Join(WorkingDir, "schema")
	)

	kent_import.NewImporter(
		dsn,
		dir,
		[]string{"id", "tenant"},
		[]string{
			"skia/pkg/ent_util",
			"entgo.io/contrib/entgql",
		},
		nil,
		rewriter,
	).Execute()
}
