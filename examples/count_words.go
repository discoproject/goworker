package main

import (
	"bufio"
	"fmt"
	"github.com/discoproject/goworker/jobutil"
	"github.com/discoproject/goworker/worker"
	"io"
	"log"
	"strings"
)

func Map(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		words := strings.Fields(text)
		for _, word := range words {
			_, err := writer.Write([]byte(word + "\n"))
			jobutil.Check(err)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("reading standard input:", err)
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
