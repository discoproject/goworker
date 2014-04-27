package worker

import (
	"fmt"
	"testing"
)

func TestInputs(t *testing.T) {
	input := []byte(`["done",[[0,"ok",0,[[0,"disco://localhost/ddfs/vol0/blob/b"]]]]]`)
	inputs := process_input(input)
	if inputs[0].replica_location != "disco://localhost/ddfs/vol0/blob/b" {
		t.Error("bad input", inputs[0])
	}
}

func TestInputsTwo(t *testing.T) {
	input := []byte(`["done",[[0,"ok",0,[[0,"disco://0"]]],[1,"ok",0,[[0,"disco://1"]]]]]`)
	inputs := process_input(input)
	if len(inputs) != 2 {
		t.Error("wrong number of inputs", len(inputs))
	}
	if inputs[0].replica_location != "disco://0" {
		t.Error("bad input", inputs[0])
	}
	if inputs[1].replica_location != "disco://1" {
		t.Error("bad input", inputs[1])
	}
}

func TestInputsMulti(t *testing.T) {
	input := []byte(`["done",[[0,"ok",0,[[0,"disco://0"]]],[1,"ok",0,[[0,"disco://1"]]],[2,"ok",0,[[0,"disco://2"]]],[3,"ok",0,[[0,"disco://3"]]],[4,"ok",0,[[0,"disco://4"]]]]]`)
	inputs := process_input(input)
	if len(inputs) != 5 {
		t.Error("wrong number of inputs", len(inputs))
	}
	for i, input := range inputs {
		if input.replica_location != "disco://"+fmt.Sprint(i) {
			t.Error("bad input", input)
		}
	}
}
