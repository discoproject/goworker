package jobutil

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var localDict map[string]string

func Setting(str string) string {
	if val := os.Getenv(str); val != "" {
		return val
	}
	if val, ok := localDict[str]; ok {
		return val
	}
	return ""
}

func SetKeyValue(key string, value string) {
	if localDict == nil {
		localDict = make(map[string]string)
	}
	localDict[key] = value
}

// TODO handle more complicated cases like when the value contains =.
func addLine(line string) {
	if line == "" {
		return
	}
	if strings.Trim(line, " \t")[0] == '#' {
		return
	}
	list := strings.Split(line, "=")
	if len(list) != 2 {
		panic("cannot process line: " + line)
	}
	key := strings.Trim(list[0], " \t\"")
	value := strings.Trim(list[1], " \t\n\"")
	SetKeyValue(key, value)
}

func addReader(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		addLine(scanner.Text())
	}
}

func AddFile(path string) {
	file, err := os.Open(path)
	Check(err)
	addReader(file)
	file.Close()
}
