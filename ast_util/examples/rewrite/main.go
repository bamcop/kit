package main

import (
	"strings"

	"github.com/bamcop/kit/ast_util/aq"
	"github.com/dave/dst"
	"github.com/helloyi/goastch/goastcher"
	"golang.org/x/exp/slog"
)

func main() {
	doc, err := aq.New("./src/src.go", nil)
	//doc, err := aq.New("", src)
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
	).DebugPrint(
		"匹配: 函数",
	).ForEach(func(root dst.Node, node dst.Node) {
		recv = node.(*dst.FuncDecl).Recv.List[0].Type.(*dst.Ident).Name
	}).FindLeave(
		goastcher.HasDescendant(goastcher.MatchCode(`entsql.Annotation{Table:`)).Bind("3"),
	).DebugPrint(
		"匹配叶子节点: entsql.Annotation",
	).Find(
		goastcher.HasDescendant(goastcher.KeyValueExpr(goastcher.Anything())).Bind("8"),
	).ForEach(func(root dst.Node, node dst.Node) {
		table = node.(*dst.KeyValueExpr).Value.(*dst.BasicLit).Value
		table = strings.Trim(table, `"`)
	}).DebugPrint("最终结果")

	slog.Info("result", slog.Any("recv", recv), slog.Any("table", table))
}
