package dyncapnp

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type SchemaMeta map[string]interface{}

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

func (m SchemaMeta) Id() uint64 {
	return getUint64(m, "id")
}

func (m SchemaMeta) Annotations() []Annotation {
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
