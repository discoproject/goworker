package jobutil

import (
	"bufio"
	"strings"
    "fmt"
	"testing"
)

func TestSorted(t *testing.T) {
	const input = "aaa\nbbb\nccc\n"
	reader := strings.NewReader(input)
	sreader := Sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)
	List := []string{"aaa", "bbb", "ccc"}
	assert_read(scanner, List, t)
}

func TestNotSorted(t *testing.T) {
	const input = "ccc\nbbb\naaa\n"
	reader := strings.NewReader(input)
	sreader := Sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)

	List := []string{"aaa", "bbb", "ccc"}
	assert_read(scanner, List, t)
}

func TestDifferntSize(t *testing.T) {
	const input = "c\nbb\naaa\n"
	reader := strings.NewReader(input)
	sreader := Sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)

	List := []string{"aaa", "bb", "c"}
	assert_read(scanner, List, t)
}

func TestUnicode(t *testing.T) {
	const input = "a\n\340\n"
	reader := strings.NewReader(input)
	sreader := Sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)

	List := []string{"a", "\340"}
	assert_read(scanner, List, t)
}

func assert_read(scanner *bufio.Scanner, List []string, t *testing.T) {
	for _, word := range List {
		scanner.Scan()
		if scanner.Text() != word {
			t.Error("not sorted", scanner.Text())
		}
	}
	if scanner.Scan() {
		t.Error("input not finished")
	}
	if err := scanner.Err(); err != nil {
		t.Error("err in input")
	}
}

func TestGrouper(t *testing.T) {
	const input = "a\na\nb\nb"
	reader := strings.NewReader(input)
	g := Grouper(reader)
	word, count, err := g.Read()
	if word != "a" {
		t.Error("wrong word", word)
	}
	if count != 2 {
		t.Error("wrong count", count)
	}
	if err != nil {
		t.Error("error in read", err)
	}

	word, count, err = g.Read()
	if word != "b" {
		t.Error("wrong word", word)
	}
	if count != 2 {
		t.Error("wrong count", count)
	}
	if err != nil {
		t.Error("error in read", err)
	}

	word, count, err = g.Read()
	if err == nil {
		t.Error("should have erred")
	}
}

func TestGrouperEmptyReader(t *testing.T) {
	const input = ""
	reader := strings.NewReader(input)
	g := Grouper(reader)
	_, _, err := g.Read()
	if err == nil {
		t.Error("should have erred")
	}
}

func TestGroupLoop(t *testing.T) {
	const input = "a\na\nb\nb"
	reader := strings.NewReader(input)
	g := Grouper(reader)

	var err error = nil
	for err == nil {
		_, count, err := g.Read()
		fmt.Println("here")
		if err != nil {
			break
		}
		if count != 2 {
			t.Error("wrong count")
		}
	}
}
