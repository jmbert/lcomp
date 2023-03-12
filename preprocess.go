package go_stack_compiler

import (
	"os/exec"
)

func preprocess(src_name string) error {
	cp_cmd := exec.Command("cp", src_name, "build/source.txt")

	err := cp_cmd.Run()

	return err
}
