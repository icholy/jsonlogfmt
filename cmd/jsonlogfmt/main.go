package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/icholy/jsonlogfmt"
)

func main() {
	var pretty bool
	var schema jsonlogfmt.Schema
	flag.Var(&schema, "field", "typed name:type fields")
	flag.BoolVar(&schema.Strict, "strict", false, "only output specified fields")
	flag.BoolVar(&pretty, "pretty", false, "pretty print the output")
	flag.Parse()
	if schema.Strict && len(schema.Fields) == 0 {
		fmt.Println("-strict requires at least one -field to be specified")
		os.Exit(1)
	}
	r := jsonlogfmt.NewReader(os.Stdin, schema)
	if pretty {
		r.SetIndent("", "  ")
	}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}
}
