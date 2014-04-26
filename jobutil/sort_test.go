package jobutil

import (
	"bufio"
	"strings"
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
	if !g.Scan() {
		t.Error("could not read")
	}
	word, count := g.Text()
	if word != "a" {
		t.Error("wrong word", word)
	}
	if count != 2 {
		t.Error("wrong count", count)
	}

	if !g.Scan() {
		t.Error("could not read")
	}
	word, count = g.Text()
	if word != "b" {
		t.Error("wrong word", word)
	}
	if count != 2 {
		t.Error("wrong count", count)
	}
	if g.Scan() {
		t.Error("reading beyond the end?")
	}
	if g.Err() != nil {
		t.Error("an error occured?")
	}
}

func TestGrouperEmptyReader(t *testing.T) {
	const input = ""
	reader := strings.NewReader(input)
	g := Grouper(reader)
	if g.Scan() {
		t.Error("should not read!")
	}
}
