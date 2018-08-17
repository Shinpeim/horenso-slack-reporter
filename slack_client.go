package reporter

import (
	"os"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackClient can post horenso result to slack
type SlackClient interface {
	Post(*horensoOut) error
}

type slackClientImpl struct {
}

type slackWebhookPayload struct {
	Attachments []slackWebhookPayloadAttachment `json:"attachments"`
}

type slackWebhookPayloadAttachment struct {
	Fallback string                               `json:"fallback"`
	Pretext  string                               `json:"pretext"`
	Color    string                               `json:"color"`
	Fields   []slackWebhookPayloadAttachmentField `json:"fields"`
}

type slackWebhookPayloadAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

func buildSlackWebhoolPayload(ho *horensoOut) *slackWebhookPayload {
	s := &slackWebhookPayload{
		Attachments: []slackWebhookPayloadAttachment{
			slackWebhookPayloadAttachment{
				Color:    "#ff0000",
				Fallback: fmt.Sprintf("command failed! %s", ho.Command),
				Pretext:  "command failed!",
				Fields: []slackWebhookPayloadAttachmentField{
					slackWebhookPayloadAttachmentField{
						Title: "command",
						Value: ho.Command,
					},
					slackWebhookPayloadAttachmentField{
						Title: "stdout",
						Value: ho.Stdout,
					},
					slackWebhookPayloadAttachmentField{
						Title: "stderr",
						Value: ho.Stderr,
					},
					slackWebhookPayloadAttachmentField{
						Title: "startAt",
						Value: ho.StartAt,
					},
					slackWebhookPayloadAttachmentField{
						Title: "endAt",
						Value: ho.EndAt,
					},
				},
			},
		},
	}

	return s
}

// NewSlackClient returns new slack client
func NewSlackClient() SlackClient {
	return &slackClientImpl{}
}

func (c *slackClientImpl) Post(ho *horensoOut) error {
	url := os.Getenv("SLACK_WEBHOOK_URL")
	if url == "" {
		return fmt.Errorf("SLACK_WEBHOOK_URL is missing")
	}

	p := buildSlackWebhoolPayload(ho)
	json, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
