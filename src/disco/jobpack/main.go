package main

import (
	"fmt"
	"os"
)

var master string
var workerDir string
var jobInputs []string

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: jobpack master_url worker_dir input(s)")
		os.Exit(1)
	}
	master = os.Args[1]
	workerDir = os.Args[2]
	jobInputs = os.Args[3:]

	CreateJobPack()
	Post(master)
	Cleanup()
}
