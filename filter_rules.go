package directus

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type FilterRule interface {
	directusFilterRule()
}

type ByField struct {
	Name   string `json:"-"`
	Filter any    `json:"-"`
}

func (ByField) directusFilterRule() {}

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

func (OR) directusFilterRule() {}

type AND struct {
	Filters []any `json:"_and"`
}

func (AND) directusFilterRule() {}

type Equal[T any] struct {
	Value T `json:"_eq"`
}

func (Equal[T]) directusFilterRule() {}

type NotEqual[T any] struct {
	Value T `json:"_neq"`
}

func (NotEqual[T]) directusFilterRule() {}

type IsOneOf[T any] struct {
	Values []T `json:"_in"`
}

func (IsOneOf[T]) directusFilterRule() {}

type IsNotOneOf[T any] struct {
	Values []T `json:"_nin"`
}

func (IsNotOneOf[T]) directusFilterRule() {}

type IsNull struct {
	Value bool `json:"_null"`
}

func (IsNull) directusFilterRule() {}

type IsNotNull struct {
	Value bool `json:"_nnull"`
}

func (IsNotNull) directusFilterRule() {}

type Contains struct {
	Value string `json:"_contains"`
}

func (Contains) directusFilterRule() {}
