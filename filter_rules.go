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
	Name   string     `json:"-"`
	Filter FilterRule `json:"-"`
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
	Filters []FilterRule `json:"_or"`
}

func (OR) directusFilterRule() {}

type AND struct {
	Filters []FilterRule `json:"_and"`
}

func (AND) directusFilterRule() {}

type Some struct {
	Filter FilterRule `json:"_some"`
}

type None struct {
	Filter FilterRule `json:"_none"`
}

type Equal[T any] struct {
	Value T `json:"_eq"`
}

func (Equal[T]) directusFilterRule() {}

type NotEqual[T any] struct {
	Value T `json:"_neq"`
}

func (NotEqual[T]) directusFilterRule() {}

type LessThan[T any] struct {
	Value T `json:"_lt"`
}

func (LessThan[T]) directusFilterRule() {}

type LessThanOrEquel[T any] struct {
	Value T `json:"_lte"`
}

func (LessThanOrEquel[T]) directusFilterRule() {}

type GreaterThan[T any] struct {
	Value T `json:"_gt"`
}

func (GreaterThan[T]) directusFilterRule() {}

type GreaterThanOrEquel[T any] struct {
	Value T `json:"_gte"`
}

func (GreaterThanOrEquel[T]) directusFilterRule() {}

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

type IsEmpty struct {
	Value bool `json:"_empty"`
}

func (IsEmpty) directusFilterRule() {}

type IsNotEmpty struct {
	Value bool `json:"_nempty"`
}

func (IsNotEmpty) directusFilterRule() {}
