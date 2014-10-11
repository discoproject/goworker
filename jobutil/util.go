package jobutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const POLL_INTERVAL = 2000

func encode(jobname string) []byte {
	list := make([]interface{}, 2)
	list[0] = POLL_INTERVAL
	list[1] = []string{jobname}
	json_jobname, err := json.Marshal(list)
	Check(err)
	return json_jobname
}

func decode_response(input []byte) (status string, results []string) {
	result := make([]interface{}, 1)
	err := json.Unmarshal(input, &result)
	Check(err)
	input0 := result[0].([]interface{})
	// jobname := input0[0].(string)

	result_list := input0[1].([]interface{})
	status = result_list[0].(string)
	inter0 := result_list[1].([]interface{})
	if len(inter0) == 0 {
		results = []string{""}
	} else {
		inter1 := inter0[0].([]interface{})
		results = make([]string, len(inter1))
		for i, item := range inter1 {
			results[i] = item.(string)
		}
	}
	return
}

func get_results(c chan []string, errChan chan error, myurl string, reqBody []byte) {
	var client *http.Client
	
	proxy := Setting("DISCO_PROXY")
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		Check(err)
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}
	resp, err := client.Post(myurl, "application/json", bytes.NewReader(reqBody))
	// TODO only retry on certain errors
	if err != nil {
		time.Sleep(time.Duration(POLL_INTERVAL) * time.Millisecond)
		get_results(c, errChan, myurl, reqBody)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		time.Sleep(time.Duration(POLL_INTERVAL) * time.Millisecond)
		get_results(c, errChan, myurl, reqBody)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errChan <- err
	}
	status, inputs := decode_response(body)
	if status == "active" {
		time.Sleep(time.Duration(POLL_INTERVAL) * time.Millisecond)
		get_results(c, errChan, myurl, reqBody)

	} else if status == "ready" {
		c <- inputs
	}
}

func Wait(master string, jobname string, timeout time.Duration) ([]string, error) {
	encoded := encode(jobname)
	url := master + "/disco/ctrl/get_results"
	c := make(chan []string)
	errChan := make(chan error)
	go get_results(c, errChan, url, encoded)
	select {
	case <-time.After(timeout * time.Second):
		return nil, errors.New("time out")
	case outputs := <-c:
		return outputs, nil
	case err := <-errChan:
		return nil, err
	}
}
