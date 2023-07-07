package aq

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/helloyi/goastch"
	"github.com/helloyi/goastch/goastcher"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

type Container struct {
	fset    *token.FileSet
	astFile *ast.File
	dstFile *dst.File
	mp      decorator.Map
	root    ast.Node
	nodes   []ast.Node
}

func New(filename string, src any) (*Container, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	dstDecorator := decorator.NewDecorator(fset)
	dstFile, err := dstDecorator.DecorateFile(astFile)
	if err != nil {
		return nil, err
	}

	container := &Container{
		fset:    fset,
		astFile: astFile,
		dstFile: dstFile,
		mp:      dstDecorator.Map,
		root:    astFile,
		nodes:   []ast.Node{astFile},
	}

	return container, nil
}

func (c *Container) Nodes() []ast.Node {
	return c.nodes
}

func (c *Container) Root() *dst.File {
	return c.dstFile
}

func (c *Container) RootString() string {
	var buff bytes.Buffer
	err := decorator.Fprint(&buff, c.dstFile)
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func (c *Container) DebugPrint(label string) *Container {
	fmt.Println(label)
	for i, node := range c.nodes {
		fmt.Println(i, strings.Repeat("-", 80))
		_ = printer.Fprint(os.Stdout, c.fset, node)
		fmt.Println(strings.Repeat("-", 80))
	}
	fmt.Println()

	return c
}

func (c *Container) DebugPrintRoot() *Container {
	fmt.Println("")

	fmt.Println("ast", strings.Repeat("#", 92))
	_ = printer.Fprint(os.Stdout, c.fset, c.root)
	fmt.Println(strings.Repeat("#", 96))

	fmt.Println("dst", strings.Repeat("#", 92))
	decorator.Print(c.dstFile)
	fmt.Println(strings.Repeat("#", 96))

	return c
}

func (c *Container) Filter(f func(node dst.Node) bool) *Container {
	c.nodes = lo.Filter(c.nodes, func(item ast.Node, index int) bool {
		var (
			dNode = c.mp.Dst.Nodes[item]
		)

		return f(dNode)
	})

	return c
}

func (c *Container) ForEach(f func(root dst.Node, node dst.Node)) *Container {
	for _, node := range c.nodes {
		node := node

		var (
			aNode = node
			dNode = c.mp.Dst.Nodes[aNode]
		)
		f(c.dstFile, dNode)
	}

	return c
}

func (c *Container) Find(rules ...goastcher.Goastcher) *Container {
	var (
		nodes []ast.Node
	)

	for _, node := range c.nodes {
		var matched []ast.Node
		for i, rule := range rules {
			results, err := goastch.Find(node, nil, rule)
			if err != nil {
				panic(err)
			}
			if len(results) == 0 {
				break
			}

			if i == 0 {
				matched = results[maps.Keys(results)[0]]
			} else {
				matched = lo.Intersect(
					matched,
					results[maps.Keys(results)[0]],
				)
			}
		}

		nodes = append(nodes, matched...)
	}

	c.nodes = nodes
	return c
}

func (c *Container) FindLeave(rules ...goastcher.Goastcher) *Container {
	var (
		nodes []ast.Node
	)

	for _, node := range c.nodes {
		var (
			matched []ast.Node
		)

		for i, rule := range rules {
			results, err := goastch.Find(node, nil, rule)
			if err != nil {
				panic(err)
			}
			if len(results) == 0 {
				matched = []ast.Node{}
				break
			}

			objs := results[maps.Keys(results)[0]]
			objs = leaveNodes(objs)
			if i == 0 {
				matched = objs
			} else {
				matched = lo.Intersect(matched, objs)
			}
		}

		nodes = append(nodes, matched...)
	}

	c.nodes = nodes
	return c
}

func (c *Container) Nth(index int, rules ...goastcher.Goastcher) *Container {
	var (
		nodes []ast.Node
	)

	for _, node := range c.nodes {
		var (
			idx     = index
			matched bool
			obj     ast.Node
		)

		for i, rule := range rules {
			// TODO: 此处假定 map 的 size == 1, 不清楚什么情况下不满足这一假设
			results, err := goastch.Find(node, nil, rule)
			if err != nil {
				panic(err)
			}
			if len(results) == 0 {
				break
			}

			objs := results[maps.Keys(results)[0]]
			if idx < 0 {
				idx = len(objs) + idx
			}
			if idx > len(objs)-1 {
				break
			}

			if i == 0 {
				obj = objs[idx]
				matched = true
			} else {
				if obj != objs[idx] {
					matched = false
					break
				}
			}
		}

		if matched == true {
			nodes = append(nodes, obj)
		}
	}

	c.nodes = nodes
	return c
}
