package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type horensoOut struct {
	Command     string   `json:"command"`
	CommandArgs []string `json:"commandArgs"`
	Output      string   `json:"output"`
	Stdout      string   `json:"stdout"`
	Stderr      string   `json:"stderr"`
	ExitCode    int      `json:"exitCode"`
	Result      string   `json:"result"`
	Pid         int      `json:"pid"`
	StartAt     string   `json:"startAt"`
	EndAt       string   `json:"endAt"`
	Hostname    string   `json:"hostName"`
	SystemTime  float32  `json:"systemTime"`
	UserTime    float32  `json:"userTime"`
}

type options struct {
	SlackWebhookURL string
	IgnoreSucceeded bool
}

func optionsFromEnv() (*options, error) {
	opts := new(options)

	opts.SlackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	if opts.SlackWebhookURL == "" {
		return nil, fmt.Errorf("SLACK_WEBHOOK_URL is missing")
	}

	doesIgnore := os.Getenv("IGNORE_SUCCEEDED")
	if doesIgnore == "" {
		opts.IgnoreSucceeded = false
	} else {
		opts.IgnoreSucceeded = true
	}

	return opts, nil
}

func parseHorensoOut(stdin io.Reader) (*horensoOut, error) {
	ho := new(horensoOut)

	text, err := ioutil.ReadAll(stdin)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(text), ho)
	return ho, err
}

// Run the reporter
func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer, c SlackClient) int {
	opts, err := optionsFromEnv()
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	ho, err := parseHorensoOut(stdin)

	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 2
	}

	if ho.ExitCode == 0 && opts.IgnoreSucceeded {
		return 0
	}

	err = c.Post(ho, opts.SlackWebhookURL)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 3
	}

	return 0
}
