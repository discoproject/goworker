package main

import (
	"archive/zip"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

const (
	MAGIC       = 0xd5c0 << 16
	MAGIC_MASK  = 0xffff << 16
	VERSION     = 0x0001
	HEADER_SIZE = 128
)

type JobPack struct {
	jobdict map[string]interface{}
	jobenv  map[string]interface{}
	jobhome string
	jobdata string
}

func (jp *JobPack) Init() {
	jp.jobdict = make(map[string]interface{})
	jp.jobenv = make(map[string]interface{})
}

func (jp *JobPack) AddToJobDict(key string, value interface{}) {
	jp.jobdict[key] = value
}

func (jp *JobPack) AddToJobEnv(key string, value interface{}) {
	jp.jobenv[key] = value
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Header struct {
	MV            uint32
	JobDictOffset uint32
	JobEnvOffset  uint32
	JobHomeOffset uint32
	JobDataOffset uint32
	_             [27]uint32
}

func compile(worker string) string {
	pwd, err := os.Getwd()
	Check(err)
    var workerDir string

	if strings.HasSuffix(worker, ".go") {
		var file string
		workerDir, file = filepath.Split(worker)
		if workerDir != "" {
			err = os.Chdir(workerDir)
			Check(err)
		}
		_, err := exec.Command("go", "build", "-o", "worker", file).Output()
		Check(err)
	} else {
        workerDir = worker
		err = os.Chdir(workerDir)
		Check(err)
		_, err := exec.Command("go", "build", "-o", "worker").Output()
		Check(err)
	}
	err = os.Chdir(pwd)
	Check(err)
	return filepath.Join(workerDir, "worker")
}

func zipit(workerExe string) string {
	//Open this executable for reading
	exeFile, err := os.Open(workerExe)
	Check(err)
	defer exeFile.Close()

	// create the zipfile
	zipfile, err := os.Create(workerExe + ".zip")
	Check(err)
	defer zipfile.Close()

	// set w to write to zipfile
	w := zip.NewWriter(zipfile)

	f, err := w.Create("job")
	Check(err)

	_, err = io.Copy(f, exeFile)
	Check(err)

	err = w.Close()
	Check(err)

	return workerExe + ".zip"
}

func Encode(jobdict map[string]interface{}, jobenv map[string]interface{},
	zipFileName string) {
	job_dict, err := json.Marshal(jobdict)
	Check(err)
	job_dict_len := len(job_dict)

	job_env, err := json.Marshal(jobenv)
	Check(err)
	job_env_len := len(job_env)

	//TODO there is no need to create the zipfile, we can actually pass the file to
	//zipit and it will zip into jp.
	zipfile, err := os.Open(zipFileName)
	Check(err)
	defer zipfile.Close()

	fileinfo, err := zipfile.Stat()
	Check(err)
	jobHomeSize := int(fileinfo.Size())

	var header Header
	header.MV = uint32(MAGIC + VERSION)
	header.JobDictOffset = uint32(HEADER_SIZE)
	header.JobEnvOffset = uint32(HEADER_SIZE + job_dict_len)
	header.JobHomeOffset = uint32(HEADER_SIZE + job_dict_len + job_env_len)
	header.JobDataOffset = uint32(HEADER_SIZE + job_dict_len + job_env_len + jobHomeSize)

	file, err := os.Create("jp")
	Check(err)
	err = binary.Write(file, binary.BigEndian, header)
	Check(err)
	binary.Write(file, binary.BigEndian, job_dict)
	binary.Write(file, binary.BigEndian, job_env)

	io.Copy(file, zipfile)
	file.Close()
}

// Decoder, currently unused
func Decode() {
	file, err := os.Open("jp")
	Check(err)
	var header Header
	err = binary.Read(file, binary.BigEndian, &header)
	// magic := header.MV & MAGIC_MASK

	job_dict_buf := make([]byte, header.JobEnvOffset-header.JobDictOffset)
	err = binary.Read(file, binary.BigEndian, &job_dict_buf)
	job_dict := make(map[string]interface{})
	err = json.Unmarshal(job_dict_buf, &job_dict)
	Check(err)

	job_env_buf := make([]byte, header.JobHomeOffset-header.JobEnvOffset)
	err = binary.Read(file, binary.BigEndian, &job_env_buf)
	job_env := make(map[string]interface{})
	err = json.Unmarshal(job_env_buf, &job_env)
	Check(err)
}

func CreateJobPack(inputs []string, worker string) {
	var jp JobPack
	jp.Init()
	host, err := os.Hostname()
	Check(err)

	user, err := user.Current()
	Check(err)

	jp.AddToJobDict("prefix", "gojob")
	jp.AddToJobDict("owner", user.Username+"@"+host)
	// Send an empty scheduler
	jp.AddToJobDict("scheduler", make(map[string]string))

	jp.AddToJobDict("reduce?", true)
	jp.AddToJobDict("save_info", "ddfs")
	jp.AddToJobDict("worker", "./job")
	jp.AddToJobDict("nr_reduces", 1)
	jp.AddToJobDict("save_results", false)

	jp.AddToJobDict("input", inputs)
	jp.AddToJobDict("map?", true)

	jp.AddToJobEnv("en", "v")
	workerExe := compile(worker)
	zipFileName := zipit(workerExe)
	Encode(jp.jobdict, jp.jobenv, zipFileName)
}

func Cleanup() {
	os.Remove("jp")                // ignore error
	os.Remove("worker/worker.zip") // ignore error

}
