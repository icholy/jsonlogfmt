package jsonlogfmt

import (
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"
)

func TestDecodeMap(t *testing.T) {
	schema := Schema{
		"number":      NumberType,
		"omg":         BoolType,
		"size":        NumberType,
		"temperature": NumberType,
	}
	f, err := os.Open(golden.Path("example.log"))
	assert.NilError(t, err)
	defer f.Close()
	r := NewReader(f, schema)
	data, err := ioutil.ReadAll(r)
	assert.NilError(t, err)
	golden.Assert(t, string(data), "example.json.golden")
}
