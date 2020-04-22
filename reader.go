package jsonlogfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-logfmt/logfmt"
)

type Reader struct {
	buf    *bytes.Buffer
	dec    *logfmt.Decoder
	enc    *json.Encoder
	schema Schema
}

func NewReader(r io.Reader, schema Schema) *Reader {
	if schema == nil {
		schema = Schema{}
	}
	var b bytes.Buffer
	return &Reader{
		buf:    &b,
		enc:    json.NewEncoder(&b),
		dec:    logfmt.NewDecoder(r),
		schema: schema,
	}
}

func (r *Reader) SetIndent(prefix, indent string) {
	r.enc.SetIndent(prefix, indent)
}

func (r *Reader) Read(data []byte) (int, error) {
	if r.buf.Len() == 0 {
		m, err := r.decodeMap()
		if err != nil {
			return 0, err
		}
		if err := r.enc.Encode(m); err != nil {
			return 0, err
		}
	}
	return r.buf.Read(data)
}

func (r *Reader) decodeMap() (map[string]interface{}, error) {
	if !r.dec.ScanRecord() {
		if r.dec.Err() == nil {
			return nil, io.EOF
		}
		return nil, r.dec.Err()
	}
	m := map[string]interface{}{}
	for r.dec.ScanKeyval() {
		key, val := string(r.dec.Key()), string(r.dec.Value())
		if typ, ok := r.schema[key]; ok {
			v, err := ParseValue(typ, val)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", key, err)
			}
			m[key] = v
		} else {
			m[key] = val
		}
	}
	return m, r.dec.Err()
}
