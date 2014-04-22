package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dinedal/goworker/jobutil"

	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func submit_job(master string) io.ReadCloser {
	file, err := os.Open("jp")
	Check(err)
	defer file.Close()

	fileinfo, err := file.Stat()
	Check(err)

	size := fileinfo.Size()
	data := make([]byte, size)
	count, err := file.Read(data)
	Check(err)
	if count != int(size) {
		panic("could not read all")
	}

	url := master + "/disco/job/new"
	resp, err := http.Post(url, "image/jpeg", bytes.NewReader(data))
	Check(err)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad response: ", resp.Status)
	}

	return resp.Body
}

func Post() {
	master = "http://" + jobutil.Setting("DISCO_MASTER") + ":" + jobutil.Setting("DISCO_PORT")
	response := submit_job(master)
	defer response.Close()
	body, err := ioutil.ReadAll(response)
	Check(err)

	result := make([]interface{}, 2)
	err = json.Unmarshal(body, &result)
	Check(err)
	jobname := result[1].(string)
	fmt.Println(jobname)
	get_results(master, jobname)
}

func get_results(master string, jobname string) {
	outputs, err := jobutil.Wait(master, jobname, 20)
	Check(err)
	for _, output := range outputs {
		fmt.Println(output)

		disco_root := jobutil.Setting("DISCO_ROOT")
		readCloser := jobutil.AddressReader(output, disco_root+"/data")
		defer readCloser.Close()
		scanner := bufio.NewScanner(readCloser)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
