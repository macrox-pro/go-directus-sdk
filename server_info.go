package directus

import (
	"context"
	"encoding/json"
)

type ServerInfo struct {
	Project struct {
		ProjectName string `json:"project_name"`
		ProjectLogo string `json:"project_logo,omitempty"`
		PublicURL   string `json:"public_url,omitempty"`
	} `json:"project"`
	Directus struct {
		Version string `json:"version"`
	} `json:"directus"`
	Node struct {
		Version string `json:"version"`
		Uptime  int64  `json:"uptime"`
	} `json:"node"`
	OS struct {
		Type     string `json:"type"`
		Version  string `json:"version"`
		Uptime   int64  `json:"uptime"`
		TotalMem int64  `json:"totalmem"`
	} `json:"os"`
}

func (c *Client) ServerInfo(ctx context.Context) (ServerInfo, error) {
	var info ServerInfo

	resp, err := c.resty.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Get("/server/info")

	if err != nil {
		return info, err
	}

	if resp.IsError() {
		var failed ErrorsPayload
		if err := json.Unmarshal(resp.Body(), &failed); err != nil {
			return info, err
		}
		return info, failed.Errors
	}

	err = json.Unmarshal(resp.Body(), &info)
	return info, err
}
