package jobutil

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func http_reader(address string) io.ReadCloser {
	resp, err := http.Get(address)
	Check(err)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("bad response: ", resp.Status)
	}
	return resp.Body
}

func absolute_disco_path(address string, disco_data string) string {
	return path.Join(disco_data, address[len("disco://"):])
}

func absolute_dir_path(address string, disco_data string) string {
	return path.Join(disco_data, address[len("dir://"):])
}

type DiscoReader struct {
	file *os.File
}

func (dr *DiscoReader) Read(p []byte) (n int, err error) {
	return dr.file.Read(p)
}

func (dr *DiscoReader) Close() error {
	return dr.file.Close()
}

func disco_reader(address string, dataDir string) io.ReadCloser {
	dr := new(DiscoReader)
	path := absolute_disco_path(address, dataDir)
	file, err := os.Open(path)
	Check(err)
	dr.file = file
	return dr
}

type DirReader struct {
	dirfile    *os.File
	scanner    *bufio.Scanner
	file       *os.File
	disco_data string
}

func (dr *DirReader) read_data(p []byte) (int, error) {
	n, err := dr.file.Read(p)
	if n == 0 && err == io.EOF {
		err = dr.file.Close()
		Check(err)
		dr.file = nil
		return dr.Read(p)
	}
	return n, err
}

func (dr *DirReader) Read(p []byte) (n int, err error) {
	if dr.file != nil {
		return dr.read_data(p)
	}
	// first read
	var line string
	if dr.scanner.Scan() {
		line = dr.scanner.Text()
	}
	if err := dr.scanner.Err(); err != nil {
		log.Fatal(err)
	}
	var address string
	var label, size int
	fmt.Sscanf(line, "%d %s %d", &label, &address, &size)
	path := absolute_disco_path(address, dr.disco_data)
	dr.file, err = os.Open(path)
	Check(err)
	return dr.read_data(p)
}

func (dr *DirReader) Close() error {
	if dr.file != nil {
		dr.file.Close()
	}
	return dr.dirfile.Close()
}

func dir_reader(address string, dataDir string) io.ReadCloser {
	dr := new(DirReader)
	path := absolute_dir_path(address, dataDir)
	file, err := os.Open(path)
	Check(err)
	dr.dirfile = file
	dr.scanner = bufio.NewScanner(dr.dirfile)
	dr.file = nil
	dr.disco_data = dataDir
	return dr
}

func AddressReader(address string, dataDir string) io.ReadCloser {
	scheme := strings.Split(address, "://")[0]
	switch scheme {
	case "http":
		fallthrough
	case "https":
		return http_reader(address)
	case "disco":
		return disco_reader(address, dataDir)
	case "dir":
		return dir_reader(address, dataDir)
	default:
		log.Fatal("Cannot read the input: ", scheme, " : ", address)
	}
	return nil
}
