package data

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bamcop/kit/debug"
	"github.com/imdario/mergo"
	"github.com/titanous/json5"
)

// MustLoadConfigJSON5 加载参数.
func MustLoadConfigJSON5[T any](cfg *T, b []byte, files ...string) {
	if err := json5.Unmarshal(b, cfg); err != nil {
		panic(err)
	}

	fName := filepath.Join(filepath.Dir(os.Args[0]), "config.json5")
	if len(files) > 0 {
		fName = files[0]
	} else if debug.MustIsRunTempDir() {
		// 如果 args[0] 在临时目录, 假定是通过 go run 的方式执行的
		fName = filepath.Join(debug.MustMainFileDir(), "config.json5")
	}

	b, err := os.ReadFile(fName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && len(files) == 0 {
			return
		} else {
			panic(err)
		}
	}

	src := new(T)
	if err := json5.Unmarshal(b, src); err != nil {
		panic(err)
	}

	if err := mergo.MergeWithOverwrite(cfg, *src); err != nil {
		panic(err)
	}
	return
}
