package js

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/bamcop/kit"
	"github.com/dop251/goja"
	"strconv"
	"strings"
)

// UnmarshalConfig expect script like `const config = {}`
// 因为 json5 不支持模版字符串, 所以使用 js 文件做为配置文件格式, 可以支持不添加转义符号的多行文字
// 因为 goja 尚不支持 es6, 否则, `export default {}` 是更好的配置文件格式
func UnmarshalConfig[T any](script string, v T) kit.Result[T] {
	ptr := new(T)
	if err := unmarshalConfig(script, ptr); err != nil {
		return kit.NewResultE[T](err)
	}
	return kit.NewResultV(*ptr)
}

func unmarshalConfig(script string, v any) error {
	vm := goja.New()

	_, err := vm.RunString(script)
	if err != nil {
		return fmt.Errorf("RunString: %w", err)
	}

	ret := vm.Get("config")
	if ret == nil {
		return fmt.Errorf("vm.Get config nil")
	}

	value, err := vm.RunString(`JSON.stringify(config, null, 2)`)
	if err != nil {
		return fmt.Errorf("JSON.stringify: %w", err)
	}

	str := value.Export().(string)
	if err := json.Unmarshal([]byte(str), v); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

// MarshalConfig 将对象导出为 `const config = {}` 的 JS 格式
// TODO: 对象的键如果非必要, 不用引号包围
// TODO: 某些包含特殊符号的键还需要处理
func MarshalConfig(v any) kit.Result[[]byte] {
	b, err := marshalConfig(v)
	if err != nil {
		return kit.NewResultE[[]byte](err)
	}

	return kit.NewResultV(b)
}

func marshalConfig(v any) ([]byte, error) {
	// 1
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json.MarshalIndent: %w", err)
	}

	// 2
	root, err := gabs.ParseJSON(b)
	if err != nil {
		return nil, fmt.Errorf("gabs.ParseJSON: %w", err)
	}

	// 3
	store := map[string]string{}
	handle(root, []string{}, store)

	// 4. 将字符串的值替换为键的 jsonpath
	for key, _ := range store {
		if _, err := root.SetP(key, key); err != nil {
			return nil, fmt.Errorf("root.SetP: %w", err)
		}
	}

	// 5.
	str := root.StringIndent("", "    ")

	// 6.
	for key, value := range store {
		oldV := `: "` + key + `"`
		newV := `: "` + value + `"`
		if len(strings.Split(value, "\n")) > 1 {
			newV = ": `" + value + "`"
		} else if strings.Contains(value, `"`) && strings.Contains(value, `'`) {
			newV = ": `" + value + "`"
		} else if strings.Contains(value, `"`) {
			newV = `: '` + value + `'`
		}

		str = strings.ReplaceAll(str, oldV, newV)
	}

	// 7
	str = "const config = " + str + "\n"

	return []byte(str), nil
}

func handle(node *gabs.Container, paths []string, store map[string]string) {
	if _, ok := node.Data().([]interface{}); ok {
		for i, container := range node.Children() {
			container := container
			handle(container, append(paths, strconv.Itoa(i)), store)
		}
	} else if _, ok := node.Data().(map[string]interface{}); ok {
		for key, container := range node.ChildrenMap() {
			container := container
			handle(container, append(paths, key), store)
		}
	} else if v, ok := node.Data().(string); ok {
		store[strings.Join(paths, ".")] = v
	} else {
		// pass
	}
}
