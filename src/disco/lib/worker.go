package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

func request_input() *Input {
	send("INPUT", "")
	_, _, line := recv()
	var mj []interface{}
	json.Unmarshal(line, &mj)

	flag := mj[0].(string)
	if flag != "done" {
		panic(flag)
	}
	_inputs := mj[1].([]interface{})
	inputs := _inputs[0].([]interface{})

	id := inputs[0].(float64)
	status := inputs[1].(string)

	label := -1
	switch t := inputs[2].(type) {
	case string:
		label = -1
	case float64:
		label = int(t)
	}
	_replicas := inputs[3].([]interface{})

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
	return input
}

func send_output(output *Output) {
	v := make([]interface{}, 3)
	v[0] = output.label
	v[1] = output.output_location //"http://example.com"
	v[2] = output.output_size

	send("OUTPUT", v)
	//TODO see if we should read the result from Disco.
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
	task   *Task
	input  *Input
	output *Output
}

type Process func(string, io.Writer, *Task)

func Run(Map Process, Reduce Process) {
	var w Worker
	send_worker()
	w.task = request_task()
	w.input = request_input()

	pwd, err := os.Getwd()
	Check(err)

	w.output = new(Output)
	if w.task.Stage == "map" {
		output_name := pwd + "/map_out"
		output, err := os.Create(output_name)
		Check(err)
		Map(w.input.replica_location, output, w.task)
		output.Close()
		w.output.output_location = "disco://" + output_name[len(w.task.Disco_data)+1:]
		output, err = os.Open(output_name)
		Check(err)
		fileinfo, err := output.Stat()
		Check(err)
		w.output.output_size = fileinfo.Size()
	} else if w.task.Stage == "map_shuffle" {
		w.output.output_location = w.input.replica_location
	} else {
		output_name := pwd + "/reduce_out"
		output, err := os.Create(output_name)
		Check(err)
		Reduce(w.input.replica_location, output, w.task)
		output.Close()
		w.output.output_location = "disco://" + output_name[len(w.task.Disco_data)+1:]
		output, err = os.Open(output_name)
		Check(err)
		fileinfo, err := output.Stat()
		Check(err)
		w.output.output_size = fileinfo.Size()
	}

	send_output(w.output)
	request_done()
}
