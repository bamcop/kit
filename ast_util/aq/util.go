package aq

import (
	"go/ast"

	"github.com/bamcop/kit/ast_util"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

func leaveNodes(nodes []ast.Node) []ast.Node {
	m := lo.SliceToMap(nodes, func(item ast.Node) (ast.Node, struct{}) {
		return item, struct{}{}
	})

	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if ast_util.IsDescendant(nodes[i], nodes[j]) {
				delete(m, nodes[i])
			}
		}
	}

	return maps.Keys(m)
}
