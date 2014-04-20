package main

import (
	"bufio"
	"disco/lib"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Map(input string, writer io.Writer, task *worker.Task) {
	reader := worker.AddressReader(input, task)
	defer reader.Close()
	body, err := ioutil.ReadAll(reader)
	Check(err)
	strBody := string(body)
	words := strings.Fields(strBody)
	for _, word := range words {
		_, err := writer.Write([]byte(word + "\n"))
		Check(err)
	}
}

func Reduce(input string, writer io.Writer, task *worker.Task) {
	reader := worker.AddressReader(input, task)
	defer reader.Close()
	//TODO sort data only if necessary
	sreader := worker.Sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)

	count := 1
	var prev, line string
	for scanner.Scan() {
		line = scanner.Text()
		if prev == "" {
			prev = line
			continue
		}
		if line == prev {
			count++
		} else {
			str := fmt.Sprintf("%d %s\n", count, prev)
			_, err := writer.Write([]byte(str))
			Check(err)
			prev = line
			count = 1
		}
	}
	str := fmt.Sprintf("%d %s\n", count, line)
	_, err := writer.Write([]byte(str))
	Check(err)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	worker.Run(Map, Reduce)
}
