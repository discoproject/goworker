package jobutil

import (
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
	if urls[0] != "disco://localhost/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff" {
		t.Error("error decoding. ", urls)
	}
}

func TestConvertDdfs(t *testing.T) {
	input := "disco://localhost/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff"
	SetKeyValue("DDFS_DATA", "/disco/ddfs/")
	path := absolute_ddfs_path(input)
	if path != "/disco/ddfs/vol0/blob/2b/train-0$574-8412a-e2ff" {
		t.Error("path not correct", path)
	}
}
