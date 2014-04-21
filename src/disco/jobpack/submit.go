package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Post(master string) {
	if !strings.HasPrefix(master, "http") {
		master = "http://" + master
	}
	url := master + "/disco/job/new"

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

	resp, err := http.Post(url, "image/jpeg", bytes.NewReader(data))
	Check(err)

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad response: ", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	Check(err)

	result := make([]interface{}, 2)
	err = json.Unmarshal(body, &result)
	Check(err)
	fmt.Println(result[1])
}
