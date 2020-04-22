package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/icholy/jsonlogfmt"
)

func main() {
	var pretty bool
	schema := jsonlogfmt.Schema{}
	flag.Var(&schema, "field", "typed name:type fields")
	flag.BoolVar(&pretty, "pretty", false, "pretty print the output")
	flag.Parse()
	r := jsonlogfmt.NewReader(os.Stdin, schema)
	if pretty {
		r.SetIndent("", "  ")
	}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}
}
