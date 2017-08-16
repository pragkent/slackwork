package sink

import (
	"bytes"
	"fmt"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pragkent/slackwork/wework"
)

var (
	SlackLinkPattern = regexp.MustCompile("<([^|>]*)\\|([^>]*)>")
)

type WeWorkSink struct {
	wc      *wework.AgentClient
	AgentID int
	tc      *TagCache
}

type TagCache struct {
	wc        *wework.AgentClient
	tags      map[string]int
	expiresAt time.Time
	mu        sync.Mutex
}

func (tc *TagCache) getTagID(name string) (int, bool) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.refresh()

	id, ok := tc.tags[name]
	return id, ok
}

func (tc *TagCache) refresh() {
	now := time.Now()
	if tc.expiresAt.After(now) {
		return
	}

	taglist, err := tc.wc.GetTagList()
	if err != nil {
		log.Printf("refresh tag list error: %v", err)
		return
	}

	for _, t := range taglist {
		tc.tags[t.Name] = t.ID
	}

	log.Printf("Tag list refreshed: %v", tc.tags)
	tc.expiresAt = now.Add(15 * time.Minute)
}

func NewWeWork(corpID string, secret string, agentID int) Sink {
	wc := wework.NewAgentClient(corpID, secret, agentID)
	return &WeWorkSink{
		wc:      wc,
		AgentID: agentID,
		tc: &TagCache{
			wc:   wc,
			tags: make(map[string]int),
		},
	}
}

func (w *WeWorkSink) Send(payload *Payload) error {
	msgs := w.Translate(payload)
	for i := range msgs {
		if err := w.wc.SendMessage(&msgs[i]); err != nil {
			return err
		}
	}

	return nil
}

func (w *WeWorkSink) Translate(payload *Payload) []wework.SendMessageRequest {
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

func (w *WeWorkSink) buildTextMessages(payload *Payload) []wework.SendMessageRequest {
	var m wework.SendMessageRequest
	m.ToTag = append(m.ToTag, w.getTagID(payload.Channel))
	m.AgentID = w.AgentID
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

	return []wework.SendMessageRequest{m}
}

func (w *WeWorkSink) buildTextCardMessages(payload *Payload) []wework.SendMessageRequest {
	var mrs []wework.SendMessageRequest

	for _, attachment := range payload.Attachments {
		var m wework.SendMessageRequest

		m.ToTag = append(m.ToTag, w.getTagID(payload.Channel))
		m.AgentID = w.AgentID
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
		fmt.Fprintf(&buf, "<div class=\"normal\">%v</div>", html.EscapeString(f.Value))
	}

	return buf.String()
}

func (w *WeWorkSink) TranslateText(text string) string {
	return w.replaceLink(text)
}

func (w *WeWorkSink) replaceLink(text string) string {
	return SlackLinkPattern.ReplaceAllString(text, "<a href=\"$1\">$2</a>")
}

func (w *WeWorkSink) getTagID(channel string) string {
	id, ok := w.tc.getTagID(strings.TrimLeft(channel, "#"))
	if !ok {
		log.Printf("tag not found for channel: %v", channel)
		return ""
	}

	return strconv.Itoa(id)
}
