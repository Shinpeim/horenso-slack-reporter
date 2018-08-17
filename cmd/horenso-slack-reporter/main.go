package main

import (
	"os"

	"github.com/Shinpeim/horenso-slack-reporter"
)

func main() {
	os.Exit(reporter.Run(os.Stdin, os.Stdout, os.Stderr, reporter.NewSlackClient()))
}
