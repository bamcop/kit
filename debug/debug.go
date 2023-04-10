package debug

import (
	"bytes"
	"debug/gosym"
	"debug/macho"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bamcop/kit"
	"github.com/lifeng1992/build_info"
)

func PrintBuildInfo() {
	info := build_info.Info()
	fmt.Printf("\n%s\n", strings.Repeat("─", 80))
	fmt.Println("GoVersion:\t", info.GoVersion)
	fmt.Println("GitCommit:\t", info.GitCommit)
	fmt.Println("BuildTime:\t", info.BuildTime)
	fmt.Printf("%s\n\n", strings.Repeat("─", 80))
}

// MainFilePath 在 init 函数中无法通过 debug.Stack() 获取 main.go 的位置
// 来源: vulncheck/internal/buildinfo/buildinfo.go: openExe
// 来源: vulncheck/internal/buildinfo/additions_scan.go: ExtractPackagesAndSymbols
// 来源: debug/gosym/pclntab_test.go: TestPCLine
func MainFilePath() kit.Result[string] {
	b, err := os.ReadFile(os.Args[0])
	if err != nil {
		return kit.NewResultE[string](err)
	}

	f, err := macho.NewFile(bytes.NewReader(b))
	if err != nil {
		return kit.NewResultE[string](err)
	}

	var textOffset uint64
	text := f.Section("__text")
	if text != nil {
		textOffset = uint64(text.Offset)
	}

	pclntab := f.Section("__gopclntab")
	if pclntab == nil {
		return kit.NewResultE[string](fmt.Errorf("gopclntab is nil"))
	}
	pclndat, err := pclntab.Data()
	if err != nil {
		return kit.NewResultE[string](err)
	}

	lineTab := gosym.NewLineTable(pclndat, textOffset)
	tab, err := gosym.NewTable(nil, lineTab)
	if err != nil {
		return kit.NewResultE[string](err)
	}

	fn := tab.LookupFunc("main.main")
	file, _, _ := tab.PCToLine(fn.Sym.Value)

	return kit.NewResultV(file)
}

func MustMainFileDir() string {
	return filepath.Dir(MainFilePath().Unwrap())
}

// MustIsRunTempDir 仅在 osx 测试过
func MustIsRunTempDir() bool {
	temp, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		panic(err)
	}

	root_path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(root_path, temp) {
		return true
	}
	return false
}
