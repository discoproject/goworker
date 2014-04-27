package jobutil

import (
	"io"
	"log"
)

type ReadClosers struct {
	rcs []io.ReadCloser
}

func (rcs *ReadClosers) init() {
	rcs.rcs = make([]io.ReadCloser, 0)
}

func (rcs *ReadClosers) add(rc io.ReadCloser) {
	rcs.rcs = append(rcs.rcs, rc)
}

func (rcs *ReadClosers) Read(p []byte) (int, error) {
	if len(rcs.rcs) == 0 {
		return 0, io.EOF
	}
	n, err := rcs.rcs[0].Read(p)
	if n != 0 {
		if err == io.EOF {
			err = nil
		}
		return n, err
	} else {
		if err == io.EOF {
			rcs.rcs[0].Close()
			rcs.rcs = rcs.rcs[1:]
		} else {
			log.Print("Got no errors and read nothing")
		}
		return rcs.Read(p)
	}
}

func (rcs *ReadClosers) Close() error {
	var err error
	for _, rc := range rcs.rcs {
		err = rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
