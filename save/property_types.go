package save

type ReadOrWritable interface {
}

type ArrayProperty struct {
	Type   string           `json:"type"`
	Values []ReadOrWritable `json:"values"`

	StructName      *string `json:"struct_name,omitempty"`
	StructType      *string `json:"struct_type,omitempty"`
	StructSize      *int    `json:"struct_size,omitempty"`
	Magic1          *[]byte `json:"magic_1,omitempty"`
	StructClassType *string `json:"struct_class_type,omitempty"`
	Magic2          *[]byte `json:"magic_2,omitempty"`
}

type StructProperty struct {
	Type  string         `json:"type"`
	Magic *[]byte        `json:"magic,omitempty"`
	Value ReadOrWritable `json:"value"`
}

type MapProperty struct {
	KeyType   string                              `json:"key_type"`
	ValueType string                              `json:"value_type"`
	Magic     []byte                              `json:"magic"`
	Values    map[string][]map[string]interface{} `json:"values"`
}

type ObjectProperty struct {
	World string `json:"world"`
	Class string `json:"class"`
}

type TextProperty struct {
	Magic  []byte `json:"magic"`
	String string `json:"string"`
}

type EnumProperty struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ByteProperty struct {
	Byte byte `json:"byte"`

	EnumType *string `json:"enum_type,omitempty"`
	EnumName *string `json:"enum_name,omitempty"`
}
