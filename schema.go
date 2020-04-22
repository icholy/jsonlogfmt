package jsonlogfmt

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Type is the type of a field
type Type int

const (
	NumberType Type = iota
	DurationType
	StringType
	TimeType
	BoolType
)

// String returns the string representation of a type
func (t Type) String() string {
	switch t {
	case NumberType:
		return "number"
	case DurationType:
		return "duration"
	case StringType:
		return "string"
	case BoolType:
		return "bool"
	default:
		return "invalid"
	}
}

// ParseType parses the string representation of a type.
// Currently: number, duration, string, bool
func ParseType(s string) (Type, error) {
	switch s {
	case "number":
		return NumberType, nil
	case "duration":
		return DurationType, nil
	case "string":
		return StringType, nil
	case "bool":
		return BoolType, nil
	default:
		return 0, fmt.Errorf("invalid type: %q", s)
	}
}

// ParseValue parses a string as the provided type
func ParseValue(t Type, s string) (interface{}, error) {
	switch t {
	case NumberType:
		v, err := strconv.ParseFloat(s, 64)
		return v, err
	case DurationType:
		d, err := time.ParseDuration(s)
		return d, err
	case BoolType:
		b, err := strconv.ParseBool(s)
		return b, err
	case StringType:
		return s, nil
	default:
		return nil, fmt.Errorf("invalid type")
	}
}

// Schema describes a set of fields and their types
type Schema struct {
	Strict bool // only output fields in the schema
	Fields map[string]Type
}

// Valid returns true if the key should be included in the output
func (s Schema) Valid(key string) bool {
	if !s.Strict {
		return true
	}
	_, ok := s.Fields[key]
	return ok
}

// Parse the key/value pair
func (s Schema) Parse(key, value string) (interface{}, error) {
	typ, ok := s.Fields[key]
	if !ok {
		return value, nil
	}
	return ParseValue(typ, value)
}

// String returns a string representation of the schema.
func (s Schema) String() string {
	var b strings.Builder
	for key, typ := range s.Fields {
		fmt.Fprintf(&b, "%s:%v\n", key, typ)
	}
	return b.String()
}

// Set implements flag.Value and allows expects fields
// with the following syntax name:type (see ParseType for valid types)
func (s *Schema) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid syntax, expecing name:type")
	}
	key := parts[0]
	typ, err := ParseType(parts[1])
	if err != nil {
		return err
	}
	if s.Fields == nil {
		s.Fields = map[string]Type{}
	}
	s.Fields[key] = typ
	return nil
}

// InferFields takes a struct type, and converts it into a set of
// fields to be used in a schema. This function will panic
// if t is not a struct type or a pointer to a struct
func InferFields(t reflect.Type) map[string]Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Errorf("expected struct, got %v", t))
	}
	fields := map[string]Type{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := fieldKey(field)
		if typ, ok := fieldType(field); ok {
			fields[name] = typ
		}
	}
	return fields
}

func fieldType(f reflect.StructField) (Type, bool) {
	t := f.Type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t {
	case reflect.TypeOf(time.Duration(0)):
		return DurationType, true
	case reflect.TypeOf(time.Time{}):
		return TimeType, true
	}
	switch t.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return NumberType, true
	case reflect.Bool:
		return BoolType, true
	case reflect.String:
		return StringType, true
	}
	return 0, false
}

func fieldKey(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" {
		return f.Name
	}
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
