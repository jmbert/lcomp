package lcomp

import (
	"fmt"
	"os"
)

type Token string

func Compile(file_name, out_name string) error {

	var file_contents string

	var buffer []byte

	err := preprocess(file_name)

	if err != nil {
		return err
	}

	buffer, err = os.ReadFile("/home/joe/lbuild/source.txt")

	if err != nil {
		return err
	}

	file_contents = string(buffer)

	tokens, err := Tokenise(file_contents)

	if err != nil {
		return err
	}

	err = Parse(tokens)

	if err != nil {
		return err
	}

	out_file, err := os.Create(out_name)

	if err != nil {
		return err
	}

	_, err = out_file.Write([]byte("extern print_char\n"))

	if err != nil {
		return err
	}

	_, err = out_file.Write([]byte("section .text\nglobal _start\n\n_start:\njmp main\n"))

	if err != nil {
		return err
	}

	for _, s := range functions {
		_, err = out_file.Write([]byte(s))

		if err != nil {
			return err
		}
	}

	_, err = out_file.Write([]byte("section .bss\n"))

	if err != nil {
		return err
	}

	for _, v := range variables {
		_, err = out_file.Write([]byte(fmt.Sprintf("%s resq 2\n", v)))

		if err != nil {
			return err
		}
	}

	return nil
}
