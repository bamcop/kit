package gin_genc

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/OneOfOne/struct2ts"
	"github.com/bamcop/kit/tool/cmd/gin_genc/helper"
	"github.com/bamcop/kit/util/hos"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"golang.org/x/tools/go/packages"
)

var (
	//go:embed assets/struct2ts.go
	struct2ts_code string
	//go:embed assets/cors.txt
	cors string
)

var (
	filter = func(path string, info fs.FileInfo) bool {
		if strings.Contains(path, "_gen") && !strings.Contains(path, "gin_genc") {
			return false
		}
		if strings.Contains(path, "tmp") {
			return false
		}
		return true
	}
)

type Option func(*AppContext)

type AppContext struct {
	root            string   // 根路径
	handler_root    string   // gin handler 根路径
	module          string   // go.mod 对应的名称
	pkgName         string   // types_gen.go 和 router_gen.go 所在的目录名称
	output          string   // types_gen.go 和 router_gen.go 所在的路径
	extraImport     []string // router_gen.go 额外需要的 import
	filter          func(fs.FileInfo) bool
	files           map[string]*ast.File
	succ            []Handler
	fail            []Handler
	axios_namespace string
	js_d_ts         string // 生成的 JavaScript d.ts 路径
	api_prefix      string // api 请求路径的前缀
	api_output      string // 生成的 JavaScript axios 请求文件路径
	ctx_provider    string // 自定义 ctx 的 New 方法
	tmp_run_dir     string // 运行临时代码的目录
}

func NewApp(srvRoot string, webRoot string, module string, options ...Option) *AppContext {
	app := &AppContext{
		root:            srvRoot,
		handler_root:    filepath.Join(srvRoot, "model"),
		module:          module,
		pkgName:         "skia_gen",
		output:          filepath.Join(srvRoot, "skia_gen"),
		extraImport:     []string{},
		filter:          nil,
		files:           nil,
		succ:            nil,
		fail:            nil,
		axios_namespace: "remote",
		js_d_ts:         filepath.Join(webRoot, "types.d.ts"),
		api_prefix:      "api",
		api_output:      filepath.Join(webRoot, "src/apis/apis.js"),
		ctx_provider:    "",
		tmp_run_dir:     filepath.Join(srvRoot, "tmp/struct2ts"),
	}

	for _, option := range options {
		option(app)
	}

	return app
}

func WithGenPkgName(pkgName string) Option {
	return func(app *AppContext) {
		app.pkgName = pkgName
		app.output = filepath.Join(app.root, pkgName)
	}
}

func WithHandlerDir(dir string) Option {
	return func(app *AppContext) {
		app.handler_root = filepath.Join(app.root, dir)
	}
}

func WithCtxProvider(ctxProvider string) Option {
	return func(app *AppContext) {
		app.ctx_provider = ctxProvider
	}
}

func WithExtraImport(extraImport []string) Option {
	return func(app *AppContext) {
		app.extraImport = extraImport
	}
}

func (app *AppContext) Start() {
	app.EnsureDir()
	app.ParserWalk()
	app.ParseHandler()
	app.MakeRouter()
	app.MakeTypes()
	app.FixImport()
	app.Tygo()
	app.WriteDTs()
	app.WriteAxios()
}

func (app *AppContext) EnsureDir() {
	hos.MustMkdirAll(filepath.Join(app.handler_root))
	hos.MustMkdirAll(filepath.Join(app.output))
	hos.MustMkdirAll(app.tmp_run_dir)
	hos.MustMkdirAll(filepath.Dir(app.api_output))
}

func (app *AppContext) ParserWalk() {
	files, err := ParserWalk(app.handler_root, filter)
	if err != nil {
		panic(err)
	}
	app.files = files
}

func (app *AppContext) ParseHandler() {
	succ, fail := ParseHandler(app.files)
	sort.Slice(succ, func(i, j int) bool {
		return succ[i].Name < succ[j].Name
	})

	app.succ = succ
	app.fail = fail

	for _, handler := range fail {
		rel, _ := filepath.Rel(app.root, handler.Filename)
		slog.Warn(handler.Error.Error(), slog.String("path", rel), slog.String("func", handler.Name))
	}

	for i, handler := range app.succ {
		rel, err := filepath.Rel(app.root, handler.Filename)
		if err != nil {
			panic(err)
		}
		dir := filepath.Dir(rel)
		app.succ[i].PkgPath = filepath.Join(app.module, dir)
	}
}

func (app *AppContext) MakeRouter() {
	var str string
	str = strings.Replace(cors, "// TODO: package", fmt.Sprintf("package %s\n", app.pkgName), -1)

	if len(app.succ) > 0 {
		var b strings.Builder
		for _, handler := range app.succ {
			b.WriteString(fmt.Sprintf("\t%s\n", strconv.Quote(handler.PkgPath)))
		}
		for _, str := range app.extraImport {
			b.WriteString(fmt.Sprintf("\t%s\n", strconv.Quote(str)))
		}
		str = strings.Replace(str, "// TODO: import", b.String(), -1)

		b.Reset()
		for _, handler := range app.succ {
			call_name := fmt.Sprintf("%s.%s", filepath.Base(handler.PkgPath), handler.Name)
			if handler.NameSpace != "" {
				call_name = fmt.Sprintf(
					"%s.%s{}.%s",
					filepath.Base(handler.PkgPath),
					handler.NameSpace,
					handler.Name,
				)
			}

			if handler.IsRawContext {
				b.WriteString(fmt.Sprintf(
					"\tr.POST(%s, ginx.Wrap(%s))\n",
					strconv.Quote(strcase.ToSnake(handler.Name)),
					call_name,
				))
			} else {
				b.WriteString(fmt.Sprintf(
					"\tr.POST(%s, ginx.WrapX(%s, %s))\n",
					strconv.Quote(strcase.ToSnake(handler.Name)),
					call_name,
					app.ctx_provider,
				))
			}
		}
		str = strings.Replace(str, "// TODO: handler", b.String(), -1)
	}

	err := ioutil.WriteFile(path.Join(app.output, "router_gen.go"), []byte(str), 0644)
	if err != nil {
		panic(err)
	}
}

func (app *AppContext) MakeTypes() {
	WriteTypes(app.pkgName, filepath.Join(app.output, "types_gen.go"), app.succ)
}

func (app *AppContext) FixImport() {
	GoFmt(app.output)
}

func (app *AppContext) Tygo() {
	var b strings.Builder
	for _, handler := range app.succ {
		b.WriteString(fmt.Sprintf("\ts.Add(%s.%sRequest{})\n", app.pkgName, handler.Name))
		b.WriteString(fmt.Sprintf("\ts.Add(%s.%sResponse{})\n", app.pkgName, handler.Name))
	}

	var str string
	str = strings.Replace(struct2ts_code, "// TODO: import", strconv.Quote(filepath.Join(app.module, app.pkgName)), -1)
	str = strings.Replace(str, "\t// TODO: s.Add", b.String(), -1)

	err := ioutil.WriteFile(filepath.Join(app.root, "tmp/struct2ts/struct2ts.go"), []byte(str), 0644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"go",
		[]string{"run", "tmp/struct2ts/struct2ts.go", "-o", app.js_d_ts}...,
	)
	cmd.Dir = app.root
	cmd.Stdout = os.Stdout

	fmt.Println("RUN", cmd.String())
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func (app *AppContext) WriteDTs() {
	b, err := ioutil.ReadFile(app.js_d_ts)
	if err != nil {
		panic(err)
	}

	s := strings.Replace(string(b), "export interface", "interface", -1)

	var buff strings.Builder
	buff.WriteString(s)
	buff.WriteString("\n")

	buff.WriteString("// Code generated by doc3. DO NOT EDIT.")
	buff.WriteString(`
// Utils: https://stackoverflow.com/a/69288824
type Expand<T> = T extends (...args: infer A) => infer R
  ? (...args: Expand<A>) => Expand<R>
  : T extends infer O
  ? { [K in keyof O]: O[K] }
  : never;

type ExpandRecursively<T> = T extends (...args: infer A) => infer R
  ? (...args: ExpandRecursively<A>) => ExpandRecursively<R>
  : T extends object
  ? T extends infer O
    ? { [K in keyof O]: ExpandRecursively<O[K]> }
    : never
  : T;
// Utils: end`)
	buff.WriteString(fmt.Sprintf("\n\ndeclare namespace %s {\n", app.axios_namespace))

	groups := lo.GroupBy(app.succ, func(item Handler) string {
		return item.NameSpace
	})
	for key, handlers := range groups {
		if key == "" {
			for _, handler := range handlers {
				buff.WriteString(fmt.Sprintf(
					"\tfunction %s(params: ExpandRecursively<%sRequest>): Promise<ExpandRecursively<%sResponse>>;\n",
					helper.Underscore(handler.Name),
					handler.Name,
					handler.Name,
				))
			}
		} else {
			buff.WriteString(fmt.Sprintf("\tdeclare namespace %s {\n", strcase.ToSnake(key)))
			for _, handler := range handlers {
				buff.WriteString(fmt.Sprintf(
					"\t\tfunction %s(params: ExpandRecursively<%sRequest>): Promise<ExpandRecursively<%sResponse>>;\n",
					helper.Underscore(handler.Name),
					handler.Name,
					handler.Name,
				))
			}
			buff.WriteString("\t}\n")
		}
	}
	buff.WriteString("}\n")

	err = ioutil.WriteFile(app.js_d_ts, []byte(buff.String()), 0644)
	if err != nil {
		panic(err)
	}
}

func (app *AppContext) WriteAxios() {
	var b strings.Builder
	b.WriteString("// Code generated by doc3. DO NOT EDIT.\n")
	b.WriteString(`import axios_instance from "./index";`)
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("class %s {\n", app.axios_namespace))

	groups := lo.GroupBy(app.succ, func(item Handler) string {
		return item.NameSpace
	})
	for key, handlers := range groups {
		if key == "" {
			for i, handler := range handlers {
				b.WriteString(fmt.Sprintf("\tstatic %s(params) {\n", helper.Underscore(handler.Name)))
				b.WriteString(fmt.Sprintf(
					"\t\treturn axios_instance.post(%s, params).then((res) => res.data);\n",
					strconv.Quote(filepath.Join(app.api_prefix, strcase.ToSnake(handler.Name))),
				))
				if i != len(handlers)-1 {
					b.WriteString("\t}\n\n")
				} else {
					b.WriteString("\t}\n")
				}
			}
		} else {
			b.WriteString(fmt.Sprintf("\tstatic %s = {\n", strcase.ToSnake(key)))
			for i, handler := range handlers {
				b.WriteString(fmt.Sprintf("\t\t%s: (params) => {\n", helper.Underscore(handler.Name)))
				b.WriteString(fmt.Sprintf(
					"\t\t\treturn axios_instance.post(%s, params).then((res) => res.data);\n",
					strconv.Quote(filepath.Join(app.api_prefix, strcase.ToSnake(handler.Name))),
				))
				if i != len(handlers)-1 {
					b.WriteString("\t\t}\n\n")
				} else {
					b.WriteString("\t\t}\n")
				}
			}
			b.WriteString("\t}\n\n")
		}
	}
	b.WriteString("}\n\n")
	b.WriteString(fmt.Sprintf("window.%s = %s\n", app.axios_namespace, app.axios_namespace))

	err := ioutil.WriteFile(app.api_output, []byte(b.String()), 0644)
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	Name         string
	NameSpace    string
	Filename     string
	PkgPath      string
	Request      *ast.StructType
	Response     *ast.StructType
	IsRawContext bool
	imports      []*ast.ImportSpec
	Error        error
}

func WriteTypes(pkgName string, path string, items []Handler) {
	var decls []ast.GenDecl
	var imports []*ast.ImportSpec
	var paths []string

	for _, h := range items {
		h := h
		decls = append(decls,
			helper.MakeASTGenDecl(h.Name+"Request", h.Request, h.PkgPath),
			helper.MakeASTGenDecl(h.Name+"Response", h.Response, h.PkgPath),
		)
		imports = append(imports, h.imports...)
		paths = append(paths, h.PkgPath)
	}

	str := helper.MakeASTFile(pkgName, decls, imports, paths)
	err := ioutil.WriteFile(path, []byte(str), 0644)
	if err != nil {
		panic(err)
	}
}

func ParseHandler(files map[string]*ast.File) ([]Handler, []Handler) {
	var succ []Handler
	var fail []Handler
	for filename, _ := range files {
		slog.Info(filename)

		// TODO: 此处会运行多次 packages.Load, 非常影响性能
		file, checker := NewTypeChecker(filename)

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if helper.TargetShape(fn) {
				v, err := helper.TargetFunc(fn, checker)
				if err != nil {
					fail = append(fail, Handler{Filename: filename, Name: fn.Name.Name, Error: err})
				} else {
					succ = append(succ, Handler{
						Name:         fn.Name.Name,
						Filename:     filename,
						Request:      v[0],
						Response:     v[1],
						IsRawContext: helper.IsRawGinContext(fn, checker),
						imports:      file.Imports,
					})
					if fn.Recv != nil {
						succ[len(succ)-1].NameSpace = fn.Recv.List[0].Type.(*ast.Ident).Name
					}
				}
			}
		}
	}
	return succ, fail
}

func ParserWalk(root string, filter func(string, fs.FileInfo) bool) (map[string]*ast.File, error) {
	pkgs := make(map[string]*ast.File)

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(info.Name()) != ".go" || strings.HasSuffix(info.Name(), "_test.go") {
			return nil
		}

		if !filter(path, info) {
			return nil
		}

		f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		pkgs[path] = f

		return nil
	})

	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

func NewTypeChecker(filename string) (*ast.File, *types.Checker) {
	cfg := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedName |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Dir: filepath.Dir(filename),
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		panic(err)
	}
	pkg := pkgs[0]

	f, ok := lo.Find(pkg.Syntax, func(item *ast.File) bool {
		return pkg.Fset.File(item.Package).Name() == filename
	})
	if !ok {
		panic(fmt.Sprintf("types.NewChecker: %s", filename))
	}

	checker := types.NewChecker(nil, pkg.Fset, pkg.Types, pkg.TypesInfo)

	return f, checker
}

func GoFmt(dir string) {
	slog.Info("goimports start...")
	defer slog.Info("goimports end")

	var buff bytes.Buffer

	cmd := exec.Command("goimports", "-v", "-w", dir)
	cmd.Stdout = &buff
	cmd.Stderr = &buff
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
