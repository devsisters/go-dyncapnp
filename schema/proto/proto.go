package proto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Proto map[string]interface{}

func getUint64(m map[string]interface{}, name string) uint64 {
	var (
		res uint64
		err error
	)
	switch n := m[name].(type) {
	case string:
		res, err = strconv.ParseUint(n, 10, 64)
	case json.Number:
		var i int64
		i, err = n.Int64()
		if err != nil {
			break
		}
		res = uint64(i)
	default:
		err = fmt.Errorf("unknown type for \"%s\": %v", name, m[name])
	}
	if err != nil {
		panic(err)
	}

	return res
}

func (m Proto) Id() uint64 {
	return getUint64(m, "id")
}

// Name to present to humans to identify this Node. You should not attempt to parse this.
// Its format could change. It is not guaranteed to be unique.
func (m Proto) DisplayName() string {
	return m["displayName"].(string)
}

// Shorter version of `DisplayName()`, chopped off by `displayNamePrefixLength`
func (m Proto) ShortDisplayName() string {
	prefix := getUint64(m, "displayNamePrefixLength")
	return m.DisplayName()[prefix:]
}

func (m Proto) Type() SchemaType {
	if _, ok := m["file"]; ok {
		return SchemaFile
	} else if _, ok := m["struct"]; ok {
		return SchemaStruct
	} else if _, ok := m["enum"]; ok {
		return SchemaEnum
	} else if _, ok := m["interface"]; ok {
		return SchemaInterface
	} else if _, ok := m["const"]; ok {
		return SchemaConst
	} else if _, ok := m["annotation"]; ok {
		return SchemaAnnotation
	} else {
		return SchemaUnknown
	}
}

type SchemaType uint16

const (
	SchemaFile SchemaType = iota
	SchemaStruct
	SchemaEnum
	SchemaInterface
	SchemaConst
	SchemaAnnotation
	SchemaUnknown
)

func (t SchemaType) String() string {
	switch t {
	case SchemaFile:
		return "file"
	case SchemaStruct:
		return "struct"
	case SchemaEnum:
		return "enum"
	case SchemaInterface:
		return "interface"
	case SchemaConst:
		return "const"
	case SchemaAnnotation:
		return "annotation"
	default:
		return ""
	}
}

func (m Proto) Annotations() []Annotation {
	anns := m["annotations"].([]interface{})
	res := make([]Annotation, len(anns))
	for i, ann := range anns {
		ann := ann.(map[string]interface{})
		res[i].Id = getUint64(ann, "id")
		res[i].Value = ann["value"].(map[string]interface{})
	}

	return res
}

type Annotation struct {
	Id    uint64
	Value Value
}

type Field map[string]interface{}

func (f Field) Name() string {
	return f["name"].(string)
}

func (f Field) CodeOrder() uint16 {
	return uint16(getUint64(f, "codeOrder"))
}

func (f Field) Annotations() []Annotation {
	anns := f["annotations"].([]interface{})
	res := make([]Annotation, len(anns))
	for i, ann := range anns {
		ann := ann.(map[string]interface{})
		res[i].Id = getUint64(ann, "id")
		res[i].Value = ann["value"].(map[string]interface{})
	}

	return res
}

func (f Field) DiscriminantValue() uint16 {
	return uint16(getUint64(f, "discriminantValue"))
}

func (f Field) Slot() (Slot, bool) {
	if s, ok := f["slot"]; ok {
		return s.(map[string]interface{}), true
	} else {
		return nil, false
	}
}

func (f Field) Group() (uint64, bool) {
	if g, ok := f["group"]; ok {
		return getUint64(g.(map[string]interface{}), "typeId"), true
	} else {
		return 0, false
	}
}

type Slot map[string]interface{}

func (s Slot) Offset() uint32 {
	return uint32(getUint64(s, "offset"))
}

func (s Slot) Type() Type {
	return s["type"].(map[string]interface{})
}

func (s Slot) DefaultValue() Value {
	return s["defaultValue"].(map[string]interface{})
}

func (s Slot) HadExplicitDefault() bool {
	return s["hadExplicitDefault"].(bool)
}

type Type map[string]interface{}