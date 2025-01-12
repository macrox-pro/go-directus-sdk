package helpers

import (
	"encoding/json"
	"net/url"
)

type URLParamJSON struct {
	Data any `json:"-" url:"-"`
}

func (p URLParamJSON) EncodeValues(key string, v *url.Values) error {
	if p.Data == nil {
		return nil
	}

	data, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}

	v.Add(key, string(data))
	return nil
}

func (p URLParamJSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Data)
}
