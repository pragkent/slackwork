package main

import "fmt"

const Version string = "0.1.0"

var GitCommit string

func printVersion() {
	fmt.Printf("slackwork version %s (%s)\n", Version, GitCommit)
}
