package sink

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pragkent/slackwork/wework"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		payload *Payload
		want    []wework.SendMessageRequest
	}{
		{
			payload: &Payload{
				Channel: "#haha",
				Parse:   "full",
				Attachments: []Attachment{
					{
						Color:    "#D63232",
						Fallback: "[Alerting] Test notification",
						Fields: []AttachmentField{
							{
								Short: true,
								Title: "High value",
								Value: "null",
							},
							{
								Short: true,
								Title: "Higher Value",
								Value: "200",
							},
							{
								Short: false,
								Title: "Error",
								Value: "This is only a test",
							},
						},
						Footer:     "Grafana v4.4.1",
						FooterIcon: "https://grafana.com/assets/img/fav32.png",
						ImageURL:   "http://grafana.org/assets/img/blog/mixed_styles.png",
						Text:       "@haha Someone is testing the alert notification within grafana.",
						Title:      "[Alerting] Test notification",
						TitleLink:  "https://grafana.com/",
					},
				},
			},
			want: []wework.SendMessageRequest{
				{
					ToUser:  nil,
					ToParty: nil,
					ToTag:   wework.RecipientSet{"#haha"},
					Type:    "textcard",
					AgentID: 12345,
					TextCard: &wework.TextCard{
						Title:       "[Alerting] Test notification",
						URL:         "https://grafana.com/",
						Description: "<div>@haha Someone is testing the alert notification within grafana.</div><div class=\"gray\">High value:</div><div>null</div><div class=\"gray\">Higher Value:</div><div>200</div><div class=\"gray\">Error:</div><div>This is only a test</div>",
					},
				},
			},
		},
		{
			payload: &Payload{
				Channel: "#haha",
				Parse:   "full",
				Attachments: []Attachment{
					{
						Color:    "#D63232",
						Fallback: "[Alerting] Test notification",
						Fields: []AttachmentField{
							{
								Short: true,
								Title: "High value",
								Value: "null",
							},
						},
						Text: "Attachment Text",
					},
				},

				Text: "Payload Text",
			},
			want: []wework.SendMessageRequest{
				{
					ToUser:  nil,
					ToParty: nil,
					ToTag:   wework.RecipientSet{"#haha"},
					Type:    "text",
					AgentID: 12345,
					Text: &wework.Text{
						Content: "Payload Text\n\nAttachment Text\nHigh value: null",
					},
				},
			},
		},
	}

	ws := &WeWorkSink{
		wc:      wework.NewAgentClient("1001", "2002", 12345),
		AgentID: 12345,
	}

	for _, tt := range tests {
		got := ws.Translate(tt.payload)
		if !cmp.Equal(got, tt.want) {
			t.Errorf("WeWorkSink.Translate error. got: %#v want: %#v", got, tt.want)
		}
	}
}
