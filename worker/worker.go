package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/discoproject/goworker/jobutil"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	DEBUG = true
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func debug(prefix string, msg interface{}) {
	if DEBUG {
		file, err := os.OpenFile("/tmp/debug", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		Check(err)
		defer file.Close()
		fmt.Fprintf(file, "%s: %v\n", prefix, msg)
	}
}

func send(key string, payload interface{}) {
	enc, err := json.Marshal(payload)
	if err != nil {
		panic("could not encode")
	}
	str := fmt.Sprintf("%s %d %s\n", key, len(enc), enc)
	fmt.Printf(str)
	debug("send", str)
}

func recv() (string, int, []byte) {
	var size int
	var status string
	fmt.Scanf("%s %d", &status, &size)
	reader := bufio.NewReader(os.Stdin)
	input := make([]byte, size)
	reader.Read(input)
	debug("recv", fmt.Sprintf("%d ", size)+string(input))
	return status, size, input
}

func send_worker() {
	type WorkerMsg struct {
		Pid     int    `json:"pid"`
		Version string `json:"version"`
	}
	wm := WorkerMsg{os.Getpid(), "1.1"}
	send("WORKER", wm)

	_, _, response := recv()
	if string(response) != "\"ok\"" {
		panic(response)
	}
}

func request_task() *Task {
	task := new(Task)
	send("TASK", "")
	_, _, line := recv()
	json.Unmarshal(line, &task)
	debug("info", task)
	return task
}

func request_input() []*Input {
	send("INPUT", "")
	_, _, line := recv()
	return process_input(line)
}

func process_input(jsonInput []byte) []*Input {
	var mj []interface{}

	json.Unmarshal(jsonInput, &mj)
	flag := mj[0].(string)
	if flag != "done" {
		// TODO support multiple passes for inputs
		panic(flag)
	}
	_inputs := mj[1].([]interface{})
	result := make([]*Input, len(_inputs))
	for index, rawInput := range _inputs {
		inputTuple := rawInput.([]interface{})

		id := inputTuple[0].(float64)
		status := inputTuple[1].(string)

		label := -1
		switch t := inputTuple[2].(type) {
		case string:
			label = -1
		case float64:
			label = int(t)
		}
		_replicas := inputTuple[3].([]interface{})

		replicas := _replicas[0].([]interface{})

		//FIXME avoid conversion to float when reading the item
		replica_id := replicas[0].(float64)
		replica_location := replicas[1].(string)

		debug("info", fmt.Sprintln(id, status, label, replica_id, replica_location))

		input := new(Input)
		input.id = int(id)
		input.status = status
		input.label = label
		input.replica_id = int(replica_id)
		input.replica_location = replica_location
		result[index] = input
	}
	return result
}

func send_output(outputs []*Output) {
	for _, output := range outputs {
		v := make([]interface{}, 3)
		v[0] = output.label
		v[1] = output.output_location //"http://example.com"
		v[2] = output.output_size

		send("OUTPUT", v)
		_, _, line := recv()
		debug("info", string(line))
	}
}

func request_done() {
	send("DONE", "")
	_, _, line := recv()
	debug("info", string(line))
}

type Task struct {
	Host       string
	Master     string
	Jobname    string
	Taskid     int
	Stage      string
	Grouping   string
	Group      string
	Disco_port int
	Put_port   int
	Disco_data string
	Ddfs_data  string
	Jobfile    string
}

type Input struct {
	id               int
	status           string
	label            int
	replica_id       int
	replica_location string
}

type Output struct {
	label           int
	output_location string
	output_size     int64
}

type Worker struct {
	task    *Task
	inputs  []*Input
	outputs []*Output
}

type Process func(io.Reader, io.Writer)

func (w *Worker) runStage(pwd string, prefix string, process Process) {
	output, err := ioutil.TempFile(pwd, prefix)
	output_name := output.Name()
	Check(err)
	defer output.Close()
	locations := make([]string, len(w.inputs))
	for i, input := range w.inputs {
		locations[i] = input.replica_location
	}

	readCloser := jobutil.AddressReader(locations, jobutil.Setting("DISCO_DATA"))
	process(readCloser, output)
	readCloser.Close()

	fileinfo, err := output.Stat()
	Check(err)

	w.outputs = make([]*Output, 1)
	w.outputs[0] = new(Output)

	absDiscoPath, err := filepath.EvalSymlinks(w.task.Disco_data)
	Check(err)
	w.outputs[0].output_location =
		"disco://" + jobutil.Setting("HOST") + "/disco/" + output_name[len(absDiscoPath)+1:]
	w.outputs[0].output_size = fileinfo.Size()
}

func Run(Map Process, Reduce Process) {
	var w Worker
	send_worker()
	w.task = request_task()

	jobutil.SetKeyValue("HOST", w.task.Host)
	master, port := jobutil.HostAndPort(w.task.Master)
	jobutil.SetKeyValue("DISCO_MASTER_HOST", master)
	if port != fmt.Sprintf("%d", w.task.Disco_port) {
		panic("port mismatch: " + port)
	}
	jobutil.SetKeyValue("DISCO_PORT", port)
	jobutil.SetKeyValue("PUT_PORT", string(w.task.Put_port))
	jobutil.SetKeyValue("DISCO_DATA", w.task.Disco_data)
	jobutil.SetKeyValue("DDFS_DATA", w.task.Ddfs_data)

	w.inputs = request_input()

	pwd, err := os.Getwd()
	Check(err)

	if w.task.Stage == "map" {
		w.runStage(pwd, "map_out_", Map)
	} else if w.task.Stage == "map_shuffle" {
		w.outputs = make([]*Output, len(w.inputs))
		for i, input := range w.inputs {
			w.outputs[i] = new(Output)
			w.outputs[i].output_location = input.replica_location
			w.outputs[i].label = input.label
			w.outputs[i].output_size = 0 // TODO find a way to calculate the size
		}
	} else {
		w.runStage(pwd, "reduce_out_", Reduce)
	}

	send_output(w.outputs)
	request_done()
}
