// this file was automatically generated using struct2ts -C -T -E -D -H -N -i -s skia/skia_gen.AFKResponse
//go:build ignore
// +build ignore

package main

import (
	"flag"
	"io"
	"log"
	"os"

	// TODO: import

	"github.com/OneOfOne/struct2ts"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var (
		out = flag.String("o", "-", "output")
		f   = os.Stdout
		err error
	)

	flag.Parse()
	if *out != "-" {
		if f, err = os.OpenFile(*out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
			panic(err)
		}
		defer f.Close()
	}
	if err = runStruct2TS(f); err != nil {
		panic(err)
	}
}

func runStruct2TS(w io.Writer) error {
	s := struct2ts.New(&struct2ts.Options{
		Indent: "	",

		NoAssignDefaults: true,
		InterfaceOnly:    true,

		NoConstructor: true,
		NoCapitalize:  false,
		MarkOptional:  false,
		NoToObject:    true,
		NoExports:     true,
		NoHelpers:     true,
		NoDate:        true,

		ES6: false,
	})

	// TODO: s.Add

	io.WriteString(w, "// this file was automatically generated, DO NOT EDIT\n")
	return s.RenderTo(w)
}
