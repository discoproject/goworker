package jobutil

import (
	"testing"
)

func TestScheme(t *testing.T) {
	input := "disco://hello"
	scheme, rest := scheme_split(input)
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
	input := "disco://localhost/05/Job@576:1c0fa:3033a/05/reduce"

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
