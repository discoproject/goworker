package main

import (
	"fmt"
	"github.com/discoproject/goworker/jobutil"
	"github.com/discoproject/goworker/worker"
	"io"
	"io/ioutil"
	"strings"
)

func Map(reader io.Reader, writer io.Writer) {
	body, err := ioutil.ReadAll(reader)
	jobutil.Check(err)
	strBody := string(body)
	words := strings.Fields(strBody)
	for _, word := range words {
		_, err := writer.Write([]byte(word + "\n"))
		jobutil.Check(err)
	}
}

func Reduce(reader io.Reader, writer io.Writer) {
	sreader := jobutil.Sorted(reader)
	grouper := jobutil.Grouper(sreader)

	for grouper.Scan() {
		word, count := grouper.Text()
		_, err := writer.Write([]byte(fmt.Sprintf("%d %s\n", count, word)))
		jobutil.Check(err)
	}
	if err := grouper.Err(); err != nil {
		jobutil.Check(err)
	}
	sreader.Close()
}

func main() {
	worker.Run(Map, Reduce)
}
