package main

import (
	"flag"
	"fmt"

	"github.com/discoproject/goworker/jobutil"

	"errors"
	"os"
	"strings"
)

type Inputs []string

func (i *Inputs) String() string {
	return fmt.Sprint(*i)
}
func (i *Inputs) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("Inputs already set")
	}
	for _, input := range strings.Split(value, ",") {
		*i = append(*i, input)
	}
	return nil
}

func main() {
	var master string
	var confFile string
	var inputs Inputs
	var worker string
	var jobtype string

	const (
		defaultMaster  = "localhost"
		masterUsage    = "The master node."
		defaultConf    = "/etc/disco/settings.py"
		confUsage      = "The setting file which contains disco settings"
		defaultWorker  = ""
		workerUsage    = "The worker directory, a .go file, or an executable"
		defaultInputs  = ""
		inputUsage     = "The comma separated list of inputs to the job."
		defaultJobType = "mapreduce"
		jobTypeUsage   = "type of the job (mapreduce or pipeline)"
	)
	flag.StringVar(&master, "Master", "", masterUsage)
	flag.StringVar(&master, "M", "", masterUsage)
	flag.StringVar(&confFile, "Conf", defaultConf, confUsage)
	flag.StringVar(&confFile, "C", defaultConf, confUsage)
	flag.StringVar(&worker, "Worker", defaultWorker, workerUsage)
	flag.StringVar(&worker, "W", defaultWorker, workerUsage)

	flag.Var(&inputs, "Inputs", inputUsage)
	flag.Var(&inputs, "I", inputUsage)

	flag.StringVar(&jobtype, "Type", defaultJobType, jobTypeUsage)
	flag.StringVar(&jobtype, "T", defaultJobType, jobTypeUsage)

	flag.Parse()

	if worker == "" || len(inputs) == 0 {
		fmt.Println("Usage: jobpack -W worker_dir -I input(s)")
		os.Exit(1)
	}

	jobutil.AddFile(confFile)
	if master != "" {
		jobutil.SetKeyValue("DISCO_MASTER_HOST", master)
	} else if jobutil.Setting("DISCO_MASTER_HOST") == "" {
		jobutil.SetKeyValue("DISCO_MASTER_HOST", defaultMaster)
	}

	CreateJobPack(inputs, worker, jobtype)
	Post()
	Cleanup()
}
