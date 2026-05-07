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

func (Some) directusFilterRule() {}

type None struct {
	Filter FilterRule `json:"_none"`
}

func (None) directusFilterRule() {}

type Has struct {
	Filter FilterRule `json:"_has"`
}

func (Has) directusFilterRule() {}

type NHas struct {
	Filter FilterRule `json:"_nhas"`
}

func (NHas) directusFilterRule() {}

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

type LessThanOrEqual[T any] struct {
	Value T `json:"_lte"`
}

func (LessThanOrEqual[T]) directusFilterRule() {}

type GreaterThan[T any] struct {
	Value T `json:"_gt"`
}

func (GreaterThan[T]) directusFilterRule() {}

type GreaterThanOrEqual[T any] struct {
	Value T `json:"_gte"`
}

func (GreaterThanOrEqual[T]) directusFilterRule() {}

type IsOneOf[T any] struct {
	Values []T `json:"_in"`
}

func (IsOneOf[T]) directusFilterRule() {}

type IsNotOneOf[T any] struct {
	Values []T `json:"_nin"`
}

func (IsNotOneOf[T]) directusFilterRule() {}

type Between[T any] struct {
	Values [2]T `json:"_between"`
}

func (Between[T]) directusFilterRule() {}

type NBetween[T any] struct {
	Values [2]T `json:"_nbetween"`
}

func (NBetween[T]) directusFilterRule() {}

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

type IContains struct {
	Value string `json:"_icontains"`
}

func (IContains) directusFilterRule() {}

type NContains struct {
	Value string `json:"_ncontains"`
}

func (NContains) directusFilterRule() {}

type StartsWith struct {
	Value string `json:"_starts_with"`
}

func (StartsWith) directusFilterRule() {}

type IStartsWith struct {
	Value string `json:"_istarts_with"`
}

func (IStartsWith) directusFilterRule() {}

type NStartsWith struct {
	Value string `json:"_nstarts_with"`
}

func (NStartsWith) directusFilterRule() {}

type EndsWith struct {
	Value string `json:"_ends_with"`
}

func (EndsWith) directusFilterRule() {}

type IEndsWith struct {
	Value string `json:"_iends_with"`
}

func (IEndsWith) directusFilterRule() {}

type NEndsWith struct {
	Value string `json:"_nends_with"`
}

func (NEndsWith) directusFilterRule() {}

type Like struct {
	Value string `json:"_like"`
}

func (Like) directusFilterRule() {}

type NLike struct {
	Value string `json:"_nlike"`
}

func (NLike) directusFilterRule() {}

type IsEmpty struct {
	Value bool `json:"_empty"`
}

func (IsEmpty) directusFilterRule() {}

type IsNotEmpty struct {
	Value bool `json:"_nempty"`
}

func (IsNotEmpty) directusFilterRule() {}
