package jobutil

import (
	"bufio"
	//"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func put_raw_data_in_file(input io.Reader) string {
	pwd, err := os.Getwd()
	Check(err)
	name := pwd + "/sorted_inputs"
	sortfile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer sortfile.Close()
	Check(err)
	// TODO take care of the size
	_, err = io.Copy(sortfile, input)
	Check(err)
	return name
}

func Sorted(input io.Reader) io.ReadCloser {
	name := put_raw_data_in_file(input)
	//TODO this version is only capable of sorting ascii files.  We need a better approach.
	err := os.Setenv("LC_ALL", "C")
	Check(err)
	out, err := exec.Command("sort", "-k", "1,1", "-T", ".", "-S", "10%", "-o", name, name).Output()
	if err != nil {
		log.Fatal("sorting input failed: ", out)
	}
	file, err := os.Open(name)
	Check(err)
	return file
}

type Group interface {
	Scan() bool
	Text() (string, int)
	Err() error
}

type group struct {
	scanner *bufio.Scanner
	line    string
	count   int
	current string
	err     error
}

func (g *group) init(input io.Reader) {
	g.scanner = bufio.NewScanner(input)
}

func (g *group) Scan() bool {
	count := 1
	var line string
	prev := g.current
	for g.scanner.Scan() {
		line = g.scanner.Text()
		if prev == "" {
			prev = line
			continue
		}
		if line == prev {
			count++
		} else {
			g.current = line
			g.line = prev
			g.count = count
			return true
		}
	}
	if g.current != "" {
		line = g.current
		g.current = ""
		g.line = line
		g.count = count
		return true
	}

	g.err = g.scanner.Err()
	return false
}

func (g *group) Text() (string, int) {
	return g.line, g.count
}

func (g *group) Err() error {
	return g.err
}

func Grouper(input io.Reader) Group {
	g := new(group)
	g.init(input)
	return g
}
