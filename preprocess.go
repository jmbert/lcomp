package lcomp

import (
	"os/exec"
)

func preprocess(src_name string) error {
	cp_cmd := exec.Command("cp", src_name, "/home/joe/lbuild/source.txt")

	err := cp_cmd.Run()

	return err
}
