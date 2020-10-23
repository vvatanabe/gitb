package main

import "fmt"

const (
	name      = "gitb"
	version   = "2.3.0"
	usage     = "A command line tool for using Backlog's git comfortably."
	usageText = name + " <command>"
)

var (
	commit string
	date   string
)

func FmtVersion() string {
	if commit == "" || date == "" {
		return version
	}
	return fmt.Sprintf("%s, build %s, date %s", version, commit, date)
}
