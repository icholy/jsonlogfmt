package jsonlogfmt

import (
	"fmt"
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
type Schema map[string]Type

// String returns a string representation of the schema.
func (s Schema) String() string {
	var b strings.Builder
	for key, typ := range s {
		fmt.Fprintf(&b, "%s:%v\n", key, typ)
	}
	return b.String()
}

// Set implements flag.Value and allows expects fields
// with the following syntax name:type (see ParseType for valid types)
func (s Schema) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid syntax, expecing name:type")
	}
	key := parts[0]
	typ, err := ParseType(parts[1])
	if err != nil {
		return err
	}
	s[key] = typ
	return nil
}
