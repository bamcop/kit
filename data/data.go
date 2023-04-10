package data

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/bamcop/kit"
	"github.com/samber/lo"
)

func MustDumpData(filename string, data any) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(filename, b, 0644); err != nil {
		panic(err)
	}
}

func MustLoadData[R any](filename string, result R) R {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	dst := new(R)
	if err := json.Unmarshal(b, dst); err != nil {
		panic(err)
	}

	return *dst
}

func MustLoadLine(filename string) []string {
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	strs := strings.Split(string(b), "\n")
	strs = lo.Map(strs, func(item string, index int) string {
		return strings.TrimSpace(item)
	})
	strs = lo.Filter(strs, func(item string, index int) bool {
		return item != ""
	})

	return strs
}

func MarshalIndent(data any) kit.Result[[]byte] {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return kit.NewResultE[[]byte](err)
	}

	return kit.NewResultV(b)
}

func Unmarshal[T any](b []byte) kit.Result[T] {
	ptr := new(T)
	if err := json.Unmarshal(b, ptr); err != nil {
		return kit.NewResultE[T](err)
	}
	return kit.NewResultV(*ptr)
}

func UnmarshalToMap(b []byte) kit.Result[map[string]any] {
	var mp map[string]interface{}
	if err := json.Unmarshal(b, &mp); err != nil {
		return kit.NewResultE[map[string]any](err)
	}
	return kit.NewResultV[map[string]any](mp)
}
