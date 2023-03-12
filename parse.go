package lcomp

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
	tokens        []Token
	position      int
	current_token Token

	asm string
}

var functions []string

var current_function *string

var variables []string

func get_arg(p *Parser) (Token, error) {
	var err error
	var arg Token
	if p.advance() {
		err = errors.New("expected argument, found EOF")
		return Token(""), err
	}
	arg = p.current_token

	if _, err_tmp := arg.parse(p); err_tmp == nil {
		err = fmt.Errorf("expected argument, found %s", string(arg))
		return Token(""), err
	}
	return arg, err
}

func is_func_call(ident string) (string, string, bool, error) {
	var match bool
	var err error

	var f_ident []rune

	var args []rune

	var open_call bool

	var close_call bool

	for _, c := range ident {

		if !open_call {
			if c == '(' {
				open_call = true
				continue
			}
			f_ident = append(f_ident, c)
		}
		if open_call && !close_call {
			if c == ')' {
				close_call = true
				continue
			}
			args = append(args, c)
		}
	}

	if open_call && close_call {
		match = true
	} else if !open_call && !close_call {
		match = false
	} else {
		err = errors.New("expected end of function call")
	}

	return string(args), string(f_ident), match, err
}

func (t *Token) parse(p *Parser) (string, error) {

	var cmd string
	var err error

	args_s, f_ident, match, err := is_func_call(string(*t))

	if err != nil {
		return "", err
	}

	if match {
		args := strings.Split(args_s, ", ")

		if args_s != "" {
			for _, arg := range args {
				cmd += fmt.Sprintf("push %s\n", arg)
			}
		}

		cmd += fmt.Sprintf("call %s\n", f_ident)

		return cmd, err
	}

	switch string(*t) {

	case "exit":
		cmd += "mov eax, 1\n"
		arg, err := get_arg(p)
		if err != nil {
			break
		}

		for _, v := range variables {
			if v == string(arg) {
				arg = "[" + Token(v) + "]"
			}
		}

		cmd += fmt.Sprintf("mov ebx, %s\n", arg)
		cmd += "int 0x80\n"
	case "return":
		arg, err := get_arg(p)
		if err != nil {
			break
		}

		for _, v := range variables {
			if v == string(arg) {
				arg = "[" + Token(v) + "]"
			}
		}

		cmd += fmt.Sprintf("mov eax, %s\n", arg)
		cmd += "ret\n"
		cmd += "pop ebp\n"
	case "define":
		arg, err := get_arg(p)
		if err != nil {
			break
		}
		var s string
		s += string(fmt.Sprintf("%s:\npush ebp\nmov ebp, esp\n", arg))
		functions = append(functions, s)
		current_function = &functions[len(functions)-1]

	case "declare":
		arg, err := get_arg(p)
		if err != nil {
			break
		}
		var_ident := arg

		variables = append(variables, string(var_ident))

	case "assign":
		arg, err := get_arg(p)
		if err != nil {
			break
		}
		var_ident := arg

		arg, err = get_arg(p)
		if err != nil {
			break
		}

		for _, v := range variables {
			if v == string(arg) {
				arg = "[" + Token(v) + "]"
			}
		}

		value := arg

		cmd += fmt.Sprintf("mov [%s], %s", var_ident, value)
	default:
		err = fmt.Errorf("invalid operation: %s", string(*t))
	}

	return cmd, err
}

func new_parser(tokens []Token) Parser {
	var p Parser

	p.tokens = tokens

	p.position = 0

	p.current_token = tokens[p.position]

	p.asm = ""

	return p
}

func (p *Parser) advance() bool {
	if p.position >= len(p.tokens)-1 {
		return true
	}
	p.position++
	p.current_token = p.tokens[p.position]
	return false
}

func Parse(tokens []Token) error {
	p := new_parser(tokens)

	var err error

	var cmd string

	for {
		cmd, err = p.current_token.parse(&p)

		if err != nil {
			break
		}

		*current_function += cmd

		if p.advance() {
			break
		}
	}

	return err
}
