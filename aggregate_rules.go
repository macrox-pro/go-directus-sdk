package directus

import (
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

type AggregateRule interface {
	query.Encoder

	directusAggregateRule()
}

type Count struct {
	Fields []string
}

func (Count) directusAggregateRule() {}

func (agg Count) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[count]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type CountDistinct struct {
	Fields []string
}

func (CountDistinct) directusAggregateRule() {}

func (agg CountDistinct) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[countDistinct]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type Sum struct {
	Fields []string
}

func (Sum) directusAggregateRule() {}

func (agg Sum) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[sum]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type SumDistinct struct {
	Fields []string
}

func (SumDistinct) directusAggregateRule() {}

func (agg SumDistinct) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[sumDistinct]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type Avg struct {
	Fields []string
}

func (Avg) directusAggregateRule() {}

func (agg Avg) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[avg]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type AvgDistinct struct {
	Fields []string
}

func (AvgDistinct) directusAggregateRule() {}

func (agg AvgDistinct) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[avgDistinct]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type Min struct {
	Fields []string
}

func (Min) directusAggregateRule() {}

func (agg Min) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[min]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type Max struct {
	Fields []string
}

func (Max) directusAggregateRule() {}

func (agg Max) EncodeValues(key string, v *url.Values) error {
	if len(agg.Fields) == 0 {
		return nil
	}

	outKey := fmt.Sprintf("%s[max]", key)
	for _, field := range agg.Fields {
		v.Add(outKey, field)
	}

	return nil
}

type Many struct {
	Rules []AggregateRule
}

func (Many) directusAggregateRule() {}

func (agg Many) EncodeValues(key string, v *url.Values) error {
	if len(agg.Rules) == 0 {
		return nil
	}

	for _, rule := range agg.Rules {
		if err := rule.EncodeValues(key, v); err != nil {
			return err
		}
	}

	return nil
}
