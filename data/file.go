package data

import (
	"errors"
	"fmt"
	"github.com/bamcop/kit"
	"os"
	"path/filepath"
)

// Lookup 从当前目录逐级向上查找指定名称的文件
func Lookup(filename string) kit.Result[string] {
	pwd, err := os.Getwd()
	if err != nil {
		return kit.NewResultE[string](fmt.Errorf("os.Getwd: %w", err))
	}

	for filepath.Base(pwd) != pwd {
		if _, err := os.Stat(filepath.Join(pwd, filename)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				pwd = filepath.Dir(pwd)
				continue
			} else {
				return kit.NewResultE[string](fmt.Errorf("os.Stat: %w", err))
			}
		}
		return kit.NewResultV(filepath.Join(pwd, filename))
	}

	return kit.NewResultE[string](os.ErrNotExist)
}
