package main

import (
	_ "embed"
	"encoding/json"
	"github.com/bamcop/kit/js"
	"os"
)

var (
	//go:embed sample_code.go.txt
	code string
	//go:embed sample_shell.sh
	shell string
	//go:embed sample_license.md
	txt string

	//go:embed input.json
	input []byte
)

type Input struct {
	Description string `json:"description"`
	Extra       struct {
		Laravel struct {
			Classmap []string `json:"classmap"`
			Code     string   `json:"code"`
		} `json:"laravel"`
	} `json:"extra"`
	Keywords []string `json:"keywords"`
	License  string   `json:"license"`
	Name     string   `json:"name"`
	Shell    string   `json:"shell"`
	Type     string   `json:"type"`
}

func main() {
	var obj Input
	if err := json.Unmarshal(input, &obj); err != nil {
		panic(err)
	}
	obj.License = txt
	obj.Shell = shell
	obj.Extra.Laravel.Code = code

	b := js.MarshalConfig(obj).Unwrap()
	if err := os.WriteFile("output.js", b, 0644); err != nil {
		panic(err)
	}
}
