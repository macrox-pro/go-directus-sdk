package helpers

import (
	"encoding"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
)

const (
	preallocateCacheSize = 128
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

type FieldsExtractor struct {
	cache map[reflect.Type][]string
	mu    sync.RWMutex
}

var DefaultFieldsExtractor = &FieldsExtractor{}

func extractFieldsByTag(tag string, prefix string, v reflect.Type) []string {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		{
			return extractFieldsByTag(tag, prefix, v.Elem())
		}
	case reflect.Struct:
		{
			const jsonTag = "json"

			if tag == jsonTag && v.Implements(jsonMarshalerType) {
				return []string{prefix}
			}

			count := v.NumField()
			if count == 0 {
				break
			}

			fields := make([]string, 0, count)
			for i := 0; i < count; i++ {
				info := v.Field(i)

				if info.Anonymous {
					if subs := extractFieldsByTag(tag, prefix, info.Type); subs != nil {
						fields = append(fields, subs...)
					}
					continue
				}

				parts := strings.SplitN(info.Tag.Get(tag), ",", 2)
				if len(parts) < 1 {
					continue
				}

				field := strings.TrimSpace(parts[0])
				if field == "-" {
					continue
				}

				if field == "" {
					field = info.Name
				}
				if prefix != "" {
					field = strings.Join([]string{prefix, field}, ".")
				}

				if subs := extractFieldsByTag(tag, field, info.Type); subs != nil {
					fields = append(fields, subs...)
				} else {
					fields = append(fields, field)
				}
			}

			if len(fields) > 0 {
				return fields
			}
		}
	}

	return nil
}

func (e *FieldsExtractor) loadFields(v reflect.Type) (fields []string, ok bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.cache == nil {
		return
	}

	fields, ok = e.cache[v]
	return
}

func (e *FieldsExtractor) setFields(v reflect.Type, fields []string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.cache == nil {
		e.cache = make(
			map[reflect.Type][]string,
			preallocateCacheSize,
		)
	}

	e.cache[v] = fields
}

func (e *FieldsExtractor) Fields(tag string, i any) []string {
	v := reflect.TypeOf(i)

	fields, ok := e.loadFields(v)
	if ok {
		return fields
	}

	fields = extractFieldsByTag(tag, "", v)
	e.setFields(v, fields)

	return fields
}

func ExtractFieldsJSON(i any) []string {
	const jsonTag = "json"

	return DefaultFieldsExtractor.Fields(jsonTag, i)
}
