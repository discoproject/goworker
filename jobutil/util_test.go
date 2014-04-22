package jobutil

import (
	"testing"
)

func TestEncodeOne(t *testing.T) {
	result := encode("gojob")
	if string(result) != "[2000,[\"gojob\"]]" {
		t.Error("bad_encoding", string(result))
	}
}

func TestDecodeActive(t *testing.T) {
	input := "[[\"gojob@576:17a3d:e4372\",[\"active\",[]]]]"
	status, result := decode_response([]byte(input))
	if status != "active" {
		t.Error("status is not correct", status)
	}

	if result[0] != "" {
		t.Error("should not find any input", result[0])
	}
}

func TestDecodeReady(t *testing.T) {
	input := "[[\"gojob@576:175ca:c4eb\",[\"ready\",[[\"disco://input\"]]]]]"
	status, result := decode_response([]byte(input))
	if status != "ready" {
		t.Error("status is not correct", status)
	}

	if result[0] != "disco://input" {
		t.Error("should not find any input", result[0])
	}
}

func SettingTest(t *testing.T) {
	SetKeyValue("hello", "world")
	val := Setting("hello")
	if val != "world" {
		t.Error("wrong value", val)
	}
}
