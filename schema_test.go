package jsonlogfmt

import (
	"reflect"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestInferFields(t *testing.T) {
	type Thing struct {
		Foo   int `json:"foo"`
		Bar   bool
		Poo   time.Duration `json:"what,"`
		Yes   string
		Buz   time.Time
		Float float64
	}
	fields := InferFields(reflect.TypeOf(&Thing{}))
	assert.DeepEqual(t, fields, map[string]Type{
		"foo":   NumberType,
		"Bar":   BoolType,
		"what":  DurationType,
		"Yes":   StringType,
		"Buz":   TimeType,
		"Float": NumberType,
	})
}
