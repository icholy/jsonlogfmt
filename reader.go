package jsonlogfmt

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/go-logfmt/logfmt"
)

// Reader converts logfmt data to json
type Reader struct {
	buf    *bytes.Buffer
	dec    *logfmt.Decoder
	enc    *json.Encoder
	schema Schema
}

// NewReader returns a Reader that converts from logfmt to json.
func NewReader(r io.Reader, schema Schema) *Reader {
	if schema.Fields == nil {
		schema.Fields = map[string]Type{}
	}
	var b bytes.Buffer
	return &Reader{
		buf:    &b,
		enc:    json.NewEncoder(&b),
		dec:    logfmt.NewDecoder(r),
		schema: schema,
	}
}

// SetIndent calls SetIndent on the json.Encoder
func (r *Reader) SetIndent(prefix, indent string) {
	r.enc.SetIndent(prefix, indent)
}

// Read implements io.Reader
func (r *Reader) Read(data []byte) (int, error) {
	if r.buf.Len() == 0 {
		m, err := r.ReadMap()
		if err != nil {
			return 0, err
		}
		if err := r.enc.Encode(m); err != nil {
			return 0, err
		}
	}
	return r.buf.Read(data)
}

// ReadMap reads the next logfmt record, converts the fields
// according to the schema and returns them in a map
func (r *Reader) ReadMap() (map[string]interface{}, error) {
	if !r.dec.ScanRecord() {
		if r.dec.Err() == nil {
			return nil, io.EOF
		}
		return nil, r.dec.Err()
	}
	m := map[string]interface{}{}
	for r.dec.ScanKeyval() {
		key, val := string(r.dec.Key()), string(r.dec.Value())
		if r.schema.Valid(key) {
			v, err := r.schema.Parse(key, val)
			if err != nil {
				return nil, err
			}
			m[key] = v
		}
	}
	return m, r.dec.Err()
}
