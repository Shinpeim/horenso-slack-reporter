package reporter

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

type mockedSlackClient struct {
	ho     *horensoOut
	Called bool
}

func (m *mockedSlackClient) Post(ho *horensoOut) error {
	m.Called = true
	return nil
}
func newMockedSlcakClient() *mockedSlackClient {
	m := &mockedSlackClient{}
	m.Called = false
	return m
}

func TestInvalidJson(t *testing.T) {
	mockedStdOut := &bytes.Buffer{}
	mockedStdErr := &bytes.Buffer{}
	invalidJSONReader := strings.NewReader("invalid json")
	mockedSlackClient := newMockedSlcakClient()

	exitCode := Run(invalidJSONReader, mockedStdOut, mockedStdErr, mockedSlackClient)

	if exitCode != 1 {
		t.Errorf("can't handle invalid JSON")
	}
	if mockedStdErr.String() == "" {
		t.Errorf("got no error message when given invalid json")
	}
}

func TestValidJson(t *testing.T) {
	os.Setenv("IGNORE_SUCCEEDED", "")
	os.Setenv("SLACK_WEBHOOK_URL", "dummy")

	json := `{
		"command": "command",
		"commandArgs": [
		  "command"
		],
		"output": "1",
		"stdout": "1",
		"stderr": "1",
		"exitCode": 0,
		"result": "command exited with code: 0",
		"pid": 95030,
		"startAt": "2015-12-28T00:37:10.494282399+09:00",
		"endAt": "2015-12-28T00:37:10.546466379+09:00",
		"hostname": "webserver.example.com",
		"systemTime": 0.034632,
		"userTime": 0.026523
	}`
	validJSONReader := strings.NewReader(json)

	_, err := parseHorensoOut(validJSONReader)

	if err != nil {
		t.Errorf("failed to parse json")
	}
}

func TestSucceededCommandIgnored(t *testing.T) {
	os.Setenv("IGNORE_SUCCEEDED", "1")
	os.Setenv("SLACK_WEBHOOK_URL", "dummy")

	json := `{
		"command": "command",
		"commandArgs": [
		  "command"
		],
		"output": "1",
		"stdout": "1",
		"stderr": "1",
		"exitCode": 0,
		"result": "command exited with code: 0",
		"pid": 95030,
		"startAt": "2015-12-28T00:37:10.494282399+09:00",
		"endAt": "2015-12-28T00:37:10.546466379+09:00",
		"hostname": "webserver.example.com",
		"systemTime": 0.034632,
		"userTime": 0.026523
	}`

	mockedStdOut := &bytes.Buffer{}
	mockedStdErr := &bytes.Buffer{}
	jr := strings.NewReader(json)
	mockedSlackClient := newMockedSlcakClient()

	exitCode := Run(jr, mockedStdOut, mockedStdErr, mockedSlackClient)

	if exitCode != 0 {
		t.Errorf("failed to handle succeeded command")
	}

	if mockedSlackClient.Called {
		t.Errorf("slack client was called when the command succeeded")
	}

}

func TestSucceededCommand(t *testing.T) {
	os.Setenv("IGNORE_SUCCEEDED", "")
	os.Setenv("SLACK_WEBHOOK_URL", "dummy")

	json := `{
		"command": "command",
		"commandArgs": [
		  "command"
		],
		"output": "1",
		"stdout": "1",
		"stderr": "1",
		"exitCode": 0,
		"result": "command exited with code: 0",
		"pid": 95030,
		"startAt": "2015-12-28T00:37:10.494282399+09:00",
		"endAt": "2015-12-28T00:37:10.546466379+09:00",
		"hostname": "webserver.example.com",
		"systemTime": 0.034632,
		"userTime": 0.026523
	}`

	mockedStdOut := &bytes.Buffer{}
	mockedStdErr := &bytes.Buffer{}
	jr := strings.NewReader(json)
	mockedSlackClient := newMockedSlcakClient()

	exitCode := Run(jr, mockedStdOut, mockedStdErr, mockedSlackClient)

	if exitCode != 0 {
		t.Errorf("failed to handle succeeded command")
	}

	if !mockedSlackClient.Called {
		t.Errorf("slack client was not called when the command succeeded")
	}

}

func TestFailedCommand(t *testing.T) {
	json := `{
		"command": "command",
		"commandArgs": [
		  "command"
		],
		"output": "1",
		"stdout": "1",
		"stderr": "1",
		"exitCode": 1,
		"result": "command exited with code: 1",
		"pid": 95030,
		"startAt": "2015-12-28T00:37:10.494282399+09:00",
		"endAt": "2015-12-28T00:37:10.546466379+09:00",
		"hostname": "webserver.example.com",
		"systemTime": 0.034632,
		"userTime": 0.026523
	}`

	mockedStdOut := &bytes.Buffer{}
	mockedStdErr := &bytes.Buffer{}
	jr := strings.NewReader(json)
	mockedSlackClient := newMockedSlcakClient()

	exitCode := Run(jr, mockedStdOut, mockedStdErr, mockedSlackClient)
	t.Log(mockedSlackClient)

	if exitCode != 0 {
		t.Errorf("failed to handle failed command")
	}

	if !mockedSlackClient.Called {
		t.Errorf("slack client wasn't called when the command failed")
	}

}
