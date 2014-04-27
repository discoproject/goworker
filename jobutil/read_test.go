package jobutil

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestScheme(t *testing.T) {
	input := "disco://hello"
	scheme, rest := SchemeSplit(input)
	if scheme != "disco" {
		t.Error("incorrect scheme", scheme)
	}

	if rest != "hello" {
		t.Error("incorrect rest", rest)
	}
}

func TestLocStr(t *testing.T) {
	input := "disco://localhost/result"
	scheme, locstr, path := loc_str(input)
	if scheme != "disco" {
		t.Error("incorrect scheme", scheme)
	}

	if locstr != "localhost" {
		t.Error("incorrect locstr", locstr)
	}

	if path != "result" {
		t.Error("incorrect path", path)
	}
}

func TestConvertUri(t *testing.T) {
	input := "disco://localhost/disco/localhost/05/Job@576:1c0fa:3033a/05/reduce"

	SetKeyValue("DISCO_MASTER", "localhost")
	SetKeyValue("DISCO_PORT", "8989")

	output := convert_uri(input)

	if output != "http://localhost:8989/disco/localhost/05/Job@576:1c0fa:3033a/05/reduce" {
		t.Error("wrong output", output)
	}
}

func TestHostPort(t *testing.T) {
	input := "http://master:8989/something"
	host, port := HostAndPort(input)
	if host != "master" {
		t.Error("wrong host", host)
	}
	if port != "8989" {
		t.Error("wrong port", port)
	}
}

func TestTagUrl(t *testing.T) {
	SetKeyValue("DISCO_MASTER", "localhost")
	SetKeyValue("DISCO_PORT", "8989")
	url := tag_url("hello")
	if url != "http://localhost:8989/ddfs/tag/hello" {
		t.Error("url not correct", url)
	}
}

func TestTagInfo(t *testing.T) {
	input := []byte(`{
               "version":1,
               "id":"train$574-8412a-18265",
               "last-modified":"2014/04/03 12:02:50",
               "urls":[["disco://localhost/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff"]],
               "user-data":{}
               }`)

	_, _, urls := tag_info(input)
	if urls[0][0] != "disco://localhost/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff" {
		t.Error("error decoding. ", urls)
	}
}

func TestConvertDdfs(t *testing.T) {
	input := "disco://localhost/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff"
	SetKeyValue("DDFS_DATA", "/disco/ddfs/")
	SetKeyValue("HOST", "localhost")
	path := absolute_ddfs_path(input)
	if path != "/disco/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff" {
		t.Error("path not correct", path)
	}
}

func TestConvertDdfsRemote(t *testing.T) {
	input := "disco://otherhost/ddfs/vol0/blob/f3/input-0$573-4ed4b-a9b1b"
	SetKeyValue("DDFS_DATA", "/disco/ddfs/")
	path := convert_uri(input)
	if path != "http://otherhost:8989/ddfs/vol0/blob/f3/input-0$573-4ed4b-a9b1b" {
		t.Error("path not correct", path)
	}
}

func TestAbsolutePath(t *testing.T) {
    input := "disco://dev02/disco/dev02/c4/gojob@576:9aa4a:ec8d/map_out_809247627"
	SetKeyValue("HOST", "dev02")
	path := absolute_disco_path(input, "/usr/local/var/disco/data/")
	if path != "/usr/local/var/disco/data/dev02/c4/gojob@576:9aa4a:ec8d/map_out_809247627" {
		t.Error("path not correct", path)
	}
}


type FakeReadCloser struct {
	reader io.Reader
}

func (frc *FakeReadCloser) Read(p []byte) (int, error) {
	return frc.reader.Read(p)
}

func (frc *FakeReadCloser) Close() error {
	return nil
}

func TestMultiReaderSingle(t *testing.T) {
	const inp1 = "aaa\n"
	var frc FakeReadCloser
	var rcs ReadClosers
	frc.reader = strings.NewReader(inp1)
	rcs.add(&frc)

	reader := bufio.NewScanner(&rcs)
	if !reader.Scan() {
		t.Error("could not read")
	}
	if text := reader.Text(); text != "aaa" {
		t.Error("wrong text: ", text)
	}
	if reader.Scan() {
		t.Error("read passed end")
	}
}

func TestMultiReader(t *testing.T) {
	var rcs ReadClosers
	const inp1 = "aaa\n"
	var frc1 FakeReadCloser
	frc1.reader = strings.NewReader(inp1)
	rcs.add(&frc1)

	var frc2 FakeReadCloser
	const inp2 = "bbb\n"
	frc2.reader = strings.NewReader(inp2)

	rcs.add(&frc2)

	reader := bufio.NewScanner(&rcs)
	if !reader.Scan() {
		t.Error("could not read first")
	}
	if text := reader.Text(); text != "aaa" {
		t.Error("wrong text: ", text)
	}
	if !reader.Scan() {
		t.Error("could not read second")
	}
	if text := reader.Text(); text != "bbb" {
		t.Error("wrong text: ", text)
	}
	if reader.Scan() {
		t.Error("read passed end")
	}
}
