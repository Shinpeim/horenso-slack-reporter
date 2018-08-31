# horenso-slack-reporter

## what's this

horenso-slack-reporter is a reporter for [horenso](https://github.com/Songmu/horenso) that post the notifications to Slack.

You must set `SLACK_WEBHOOK_URL` as your slack incomming webhook url to run the command.

And you can ignore succeeded (exited with 0) command by setting enviroment variable `IGNORE_SUCCEEDED=1` (optional).

## install

go get github.com/Shinpeim/horenso-slack-reporter

## usage

```
$ horenso --reporter /path/to/horenso-slack-reporter -- /path/to/yourjob
```