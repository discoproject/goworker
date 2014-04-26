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

    var err error = nil
	for err == nil {
        word, count, err := grouper.Read()
        if err != nil {
            break
        }
		_, err = writer.Write([]byte(fmt.Sprintf("%d %s\n", count, word)))
		jobutil.Check(err)
	}
	sreader.Close()
}

func main() {
	worker.Run(Map, Reduce)
}
