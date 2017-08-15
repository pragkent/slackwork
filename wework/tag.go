package wework

import "fmt"

const (
	GetTagListURL = "https://qyapi.weixin.qq.com/cgi-bin/tag/list"
)

type GetTagListResponse struct {
	Code    int    `json:"errcode,omitempty"`
	Message string `json:"errmsg,omitempty"`
	TagList []Tag  `json:"taglist,omitempty"`
}

type Tag struct {
	ID   int    `json:"tagid,omitempty"`
	Name string `json:"tagname,omitempty"`
}

func (c *AgentClient) GetTagList() ([]Tag, error) {
	var resp GetTagListResponse

	if err := c.getJSON(GetTagListURL, &resp); err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("tag.list api error: %+v", resp)
	}

	return resp.TagList, nil
}
