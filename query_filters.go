package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ByField struct {
	Name   string `json:"-"`
	Filter any    `json:"-"`
}

func (f ByField) MarshalJSON() ([]byte, error) {
	buf := bytes.Buffer{}

	sub, err := json.Marshal(f.Filter)
	if err != nil {
		return nil, err
	}

	if _, err = buf.WriteString(
		fmt.Sprintf(`{"%s":`, f.Name),
	); err != nil {
		return nil, err
	}

	if _, err = buf.Write(sub); err != nil {
		return nil, err
	}

	if err = buf.WriteByte('}'); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type OR struct {
	Filters []any `json:"_or"`
}

type AND struct {
	Filters []any `json:"_and"`
}

type Equal[T any] struct {
	Value T `json:"_eq"`
}

type NotEqual[T any] struct {
	Value T `json:"_neq"`
}

type IsOneOf[T any] struct {
	Values []T `json:"_in"`
}

type IsNotOneOf[T any] struct {
	Values []T `json:"_nin"`
}

type IsNull struct {
	Value bool `json:"_null"`
}

type IsNotNull struct {
	Value bool `json:"_nnull"`
}

type Contains struct {
	Value string `json:"_contains"`
}
