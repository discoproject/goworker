package main

import (
	"os"
)

var master string
var workerDir string
var jobInputs []string

func main() {
	master = os.Args[1]
	workerDir = os.Args[2]
	jobInputs = os.Args[3:]

	CreateJobPack()
	Post(master)
	Cleanup()
}
