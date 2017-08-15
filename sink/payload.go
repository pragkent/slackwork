package sink

import (
	"bytes"
	"encoding/json"
)

type Payload struct {
	Username    string       `json:"username,omitempty"`
	Text        string       `json:"text,omitempty"`
	Fallback    string       `json:"fallback,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Parse       string       `json:"parse,omitempty"`
}

type Attachment struct {
	Color      string            `json:"color,omitempty"`
	Title      string            `json:"title,omitempty"`
	TitleLink  string            `json:"title_link,omitempty"`
	Fallback   string            `json:"fallback,omitempty"`
	Text       string            `json:"text,omitempty"`
	ImageURL   string            `json:"image_url,omitempty"`
	Footer     string            `json:"footer,omitempty"`
	FooterIcon string            `json:"footer_icon,omitempty"`
	Fields     []AttachmentField `json:"fields,omitempty"`
}

type AttachmentField struct {
	Short bool   `json:"short,omitempty"`
	Value string `json:"value,omitempty"`
	Title string `json:"title,omitempty"`
}

func (f *AttachmentField) UnmarshalJSON(b []byte) error {
	var af attachmentField
	if err := json.Unmarshal(b, &af); err != nil {
		return err
	}

	value, err := af.ValueString()
	if err != nil {
		return err
	}

	f.Short = af.Short
	f.Value = value
	f.Title = af.Title

	return nil
}

type attachmentField struct {
	Short bool            `json:"short,omitempty"`
	Value json.RawMessage `json:"value,omitempty"`
	Title string          `json:"title,omitempty"`
}

func (af *attachmentField) ValueString() (string, error) {
	var value string

	if af.valueIsQuoted() {
		if err := json.Unmarshal(af.Value, &value); err != nil {
			return "", err
		}
	} else {
		value = string(af.Value)
	}

	return value, nil
}

func (af *attachmentField) valueIsQuoted() bool {
	return bytes.IndexRune(af.Value, '"') == 0
}
