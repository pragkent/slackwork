package wework

import (
	"fmt"
)

const (
	SendChatMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send"
)

type MessageType string

var (
	MessageTypeText     MessageType = "text"
	MessageTypeTextCard MessageType = "textcard"
	MessageTypeImage    MessageType = "image"
)

type Text struct {
	Content string `json:"content,omitempty"`
}

type TextCard struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	BtnTxt      string `json:"btntxt,omitempty"`
}

type Image struct {
	MediaID string `json:"media_id,omitempty"`
}

type SendChatMessageRequest struct {
	ChatID   string      `json:"chatid,omitempty"`
	Type     MessageType `json:"msgtype,omitempty"`
	Safe     int         `json:"safe,omitempty"`
	Text     *Text       `json:"text,omitempty"`
	TextCard *TextCard   `json:"textcard,omitempty"`
	Image    *Image      `json:"image,omitempty"`
}

type SendChatMessageResponse struct {
	Code    int    `json:"errcode,omitempty"`
	Message string `json:"errmsg,omitempty"`
}

func (c *AgentClient) SendChatMessage(req *SendChatMessageRequest) error {
	var resp SendChatMessageResponse
	if err := c.postJSON(SendChatMessageURL, req, &resp); err != nil {
		return err
	}

	if resp.Code != 0 {
		return fmt.Errorf("message.send api error: %+v", resp)
	}

	return nil
}
