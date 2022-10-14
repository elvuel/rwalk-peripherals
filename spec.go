package rwalk_peripherals

import (
	"bytes"
	"encoding/json"
	"io"
)

type InterfaceMapping map[string][]string

type SchemaSpec struct {
	InterfaceMapping InterfaceMapping `json:"implementsMapping"`
	UniquePkgPaths   []string
	PkgAliasMapping  map[string]string
	ParsedStructs    map[string]*StructSpec `json:"parsedStructs"`
}

func (spec *SchemaSpec) BuildUniquePkgPaths() {
	hitted := make(map[string]int8)
	result := make([]string, 0)
	for _, st := range spec.ParsedStructs {
		if st.Schema == nil {
			st.Schema = spec
		}
		if _, ok := hitted[st.Pkgpath]; !ok {
			hitted[st.Pkgpath] = 1
			result = append(result, st.Pkgpath)
		}
	}
	spec.UniquePkgPaths = result
}

func (spec *SchemaSpec) AddPkgPathAlias(pkg, alias string) {
	spec.PkgAliasMapping[pkg] = alias
}

func NewSchemaWithReader(r io.Reader) (*SchemaSpec, error) {
	schema := &SchemaSpec{}
	decoder := json.NewDecoder(r)
	err := decoder.Decode(schema)
	if err != nil {
		return nil, err
	}

	schema.BuildUniquePkgPaths()

	if schema.PkgAliasMapping == nil {
		schema.PkgAliasMapping = make(map[string]string)
	}
	return schema, nil
}

func NewSchema(data []byte) (*SchemaSpec, error) {
	return NewSchemaWithReader(bytes.NewBuffer(data))
}

type StructSpec struct {
	Kind    string             `json:"kind"`
	Pkgpath string             `json:"pkgpath"`
	Name    string             `json:"typeName"`
	Fields  []*StructFieldSpec `json:"fields"`

	ExtraData map[string]interface{} `json:"E,omitempty"`

	Schema *SchemaSpec `json:"-"`
}

func (spec *StructSpec) Clone() *StructSpec {
	newSpec := &StructSpec{}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(spec)
	json.NewDecoder(&buf).Decode(newSpec)
	newSpec.Schema = spec.Schema
	return newSpec
}

func (spec *StructSpec) AliasPkgPath() string {
	if spec.Schema != nil {
		if v, ok := spec.Schema.PkgAliasMapping[spec.Pkgpath]; ok {
			return v
		} else {
			return spec.Pkgpath
		}
	}
	return spec.Pkgpath
}

type StructFieldSpec struct {
	Name                string `json:"name,omitempty"`
	Tag                 string `json:"tag,omitempty"`
	IsPtr               string `json:"isPtr,omitempty"`
	Pkgpath             string `json:"pkgpath,omitempty"`
	TypeName            string `json:"typeName,omitempty"`
	Kind                string `json:"kind,omitempty"`
	SliceOf             string `json:"sliceOf,omitempty"`
	IsInterface         string `json:"isInterface,omitempty"`
	ArraySize           string `json:"arraySize,omitempty"`
	MapValueKind        string `json:"mapValueKind,omitempty"`
	MapValueTypeName    string `json:"mapValueTypeName,omitempty"`
	MapValueTypePkgpath string `json:"mapValueTypePkgpath,omitempty"`
	MapKeyKind          string `json:"mapKeyKind,omitempty"`
	MapKeyTypeName      string `json:"mapKeyTypeName,omitempty"`
	MapKeyTypePkgpath   string `json:"mapKeyTypePkgpath,omitempty"`

	ExtraData map[string]interface{} `json:"E,omitempty"`
}
