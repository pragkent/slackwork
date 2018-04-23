package sink

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/pragkent/slackwork/wework"
)

var (
	SlackLinkPattern = regexp.MustCompile("<([^|>]*)\\|([^>]*)>")
)

type WeWorkSink struct {
	wc      *wework.AgentClient
	AgentID int
}

func NewWeWork(corpID string, secret string, agentID int) Sink {
	wc := wework.NewAgentClient(corpID, secret, agentID)
	return &WeWorkSink{
		wc:      wc,
		AgentID: agentID,
	}
}

func (w *WeWorkSink) Send(payload *Payload) error {
	msgs := w.Translate(payload)
	for i := range msgs {
		if err := w.wc.SendChatMessage(&msgs[i]); err != nil {
			return err
		}
	}

	return nil
}

func (w *WeWorkSink) Translate(payload *Payload) []wework.SendChatMessageRequest {
	if w.shouldUseTextCard(payload) {
		return w.buildTextCardMessages(payload)
	} else {
		return w.buildTextMessages(payload)
	}
}

func (w *WeWorkSink) shouldUseTextCard(payload *Payload) bool {
	if len(payload.Text) != 0 {
		return false
	}

	if len(payload.Attachments) == 0 {
		return false
	}

	for _, a := range payload.Attachments {
		if a.TitleLink == "" {
			return false
		}
	}

	return true
}

func (w *WeWorkSink) buildTextMessages(payload *Payload) []wework.SendChatMessageRequest {
	var m wework.SendChatMessageRequest
	m.ChatID = w.getChatID(payload.Channel)
	m.Type = wework.MessageTypeText

	var buf bytes.Buffer
	buf.WriteString(w.TranslateText(payload.Text))

	for _, a := range payload.Attachments {
		buf.WriteString("\n\n")

		if a.TitleLink != "" {
			fmt.Fprintf(&buf, "<a href=\"%s\">%s</a>\n", a.TitleLink, a.Title)
		}

		if a.Text != "" {
			buf.WriteString(w.TranslateText(a.Text))
		}

		for _, f := range a.Fields {
			buf.WriteString("\n")
			fmt.Fprintf(&buf, "%s: %s", f.Title, f.Value)
		}
	}

	m.Text = &wework.Text{
		Content: buf.String(),
	}

	return []wework.SendChatMessageRequest{m}
}

func (w *WeWorkSink) buildTextCardMessages(payload *Payload) []wework.SendChatMessageRequest {
	var mrs []wework.SendChatMessageRequest

	for _, attachment := range payload.Attachments {
		var m wework.SendChatMessageRequest

		m.ChatID = w.getChatID(payload.Channel)
		m.Type = wework.MessageTypeTextCard
		m.TextCard = &wework.TextCard{
			Title:       attachment.Title,
			URL:         attachment.TitleLink,
			Description: w.getDescription(&attachment),
		}

		mrs = append(mrs, m)
	}

	return mrs
}

func (w *WeWorkSink) getDescription(a *Attachment) string {
	var buf bytes.Buffer

	if a.Text != "" {
		fmt.Fprintf(&buf, "<div class=\"normal\">%s\n</div>", html.EscapeString(a.Text))
	}

	for _, f := range a.Fields {
		fmt.Fprintf(&buf, "<div class=\"gray\">%s</div>", html.EscapeString(f.Title))
		fmt.Fprintf(&buf, "<div class=\"highlight\">%v</div>", html.EscapeString(f.Value))
	}

	if text := buf.String(); text != "" {
		return text
	} else {
		return a.Title
	}
}

func (w *WeWorkSink) TranslateText(text string) string {
	return w.replaceLink(text)
}

func (w *WeWorkSink) replaceLink(text string) string {
	return SlackLinkPattern.ReplaceAllString(text, "<a href=\"$1\">$2</a>")
}

func (w *WeWorkSink) getChatID(channel string) string {
	return strings.TrimLeft(channel, "#")
}
