package worker

import (
	"bufio"
	"strings"
	"testing"
)

func TestSorted(t *testing.T) {
	const input = "aaa\nbbb\nccc\n"
	reader := strings.NewReader(input)
	sreader := sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)
	List := []string{"aaa", "bbb", "ccc"}
	assert_read(scanner, List, t)
}

func TestNotSorted(t *testing.T) {
	const input = "ccc\nbbb\naaa\n"
	reader := strings.NewReader(input)
	sreader := sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)

	List := []string{"aaa", "bbb", "ccc"}
	assert_read(scanner, List, t)
}

func TestDifferntSize(t *testing.T) {
	const input = "c\nbb\naaa\n"
	reader := strings.NewReader(input)
	sreader := sorted(reader)
	defer sreader.Close()
	scanner := bufio.NewScanner(sreader)

	List := []string{"aaa", "bb", "c"}
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
