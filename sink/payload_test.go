package sink

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmashalPayload(t *testing.T) {
	tests := []struct {
		text string
		want *Payload
	}{
		{
			text: `{
				"attachments":[
				{
					"color":"#D63232",
					"fallback":"[Alerting] Test notification",
					"fields":[
					{
						"short":true,
						"title":"High value",
						"value":null
					},
					{
						"short":true,
						"title":"Higher Value",
						"value":200
					},
					{
						"short":false,
						"title":"Error",
						"value":"This is only a test"
					}
					],
					"footer":"Grafana v4.4.1",
					"footer_icon":"https://grafana.com/assets/img/fav32.png",
					"image_url":"http://grafana.org/assets/img/blog/mixed_styles.png",
					"text":"@haha Someone is testing the alert notification within grafana.",
					"title":"[Alerting] Test notification",
					"title_link":"https://grafana.com/"
				}
				],
				"channel":"#haha",
				"parse":"full"
			}
			`,
			want: &Payload{
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
		},
	}

	for _, tt := range tests {
		var pl Payload
		if err := json.Unmarshal([]byte(tt.text), &pl); err != nil {
			t.Errorf("json.Unmarshal error: %v", err)
		}

		if !cmp.Equal(&pl, tt.want) {
			t.Errorf("json.Unmarshal not match. got: %+v; want: %+v", &pl, tt.want)
		}
	}
}
