package ast_util

import (
	"go/ast"
)

func IsDescendant(root ast.Node, child ast.Node) bool {
	if root == child {
		return false
	}

	var (
		result = false
	)
	ast.Inspect(root, func(node ast.Node) bool {
		if node == child {
			result = true
		}
		return true
	})

	return result
}
