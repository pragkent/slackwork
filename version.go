package main

import "fmt"

const Version string = "1.0.0"

var GitCommit string

func printVersion() {
	fmt.Printf("slackwork version %s (%s)\n", Version, GitCommit)
}
