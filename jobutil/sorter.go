package jobutil

import (
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
