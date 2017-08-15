package wework

import (
	"fmt"
	"strings"
)

const (
	SendMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

type MessageType string

var (
	MessageTypeText     MessageType = "text"
	MessageTypeTextCard MessageType = "textcard"
)

type RecipientSet []string

func (f RecipientSet) MarshalJSON() ([]byte, error) {
	return []byte(strings.Join(f, "|")), nil
}

type SendMessageRequest struct {
	ToUser   RecipientSet `json:"touser,omitempty"`
	ToParty  RecipientSet `json:"toparty,omitempty"`
	ToTag    RecipientSet `json:"totag,omitempty"`
	Type     MessageType  `json:"msgtype,omitempty"`
	AgentID  int          `json:"agentid,omitempty"`
	Safe     int          `json:"safe,omitempty"`
	Text     *Text        `json:"text,omitempty"`
	TextCard *TextCard    `json:"textcard,omitempty"`
}

type SendMessageResponse struct {
	Code         int    `json:"errcode,omitempty"`
	Message      string `json:"errmsg,omitempty"`
	InvalidUser  string `json:"invaliduser,omitempty"`
	InvalidParty string `json:"invalidparty,omitempty"`
	InvalidTag   string `json:"invalidtag,omitempty"`
}

type Text struct {
	Content string `json:"content,omitempty"`
}

type TextCard struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	BtnTxt      string `json:"btntxt,omitempty"`
}

func (c *AgentClient) SendMessage(req *SendMessageRequest) error {
	var resp SendMessageResponse
	if err := c.postJSON(SendMessageURL, req, &resp); err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("message.send api error: %+v", resp)
	}

	return nil
}
