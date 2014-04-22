package main

import (
	"fmt"

	"github.com/dinedal/goworker/jobutil"

	"os"
)

var master string
var workerDir string
var jobInputs []string

// TODO add options instead of using the positional arguments
func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: jobpack discoConfFile master_url worker_dir input(s)")
		os.Exit(1)
	}
	jobutil.SetKeyValue("DISCO_MASTER", os.Args[1])
	workerDir = os.Args[2]
	confFile := os.Args[3]
	jobutil.AddFile(confFile)
	jobInputs = os.Args[4:]

	CreateJobPack()
	Post()
	Cleanup()
}
