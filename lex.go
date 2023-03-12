package lcomp

import "fmt"

type Lexer struct {
	src          string
	position     int
	current_char byte
	tokens       []Token
}

func new_lexer(src string) Lexer {
	var l Lexer
	l.src = src
	l.position = -1
	l.advance()
	return l
}

func (l *Lexer) advance() bool {
	if l.position >= len(l.src)-1 {
		return true
	}
	l.position++
	l.current_char = l.src[l.position]
	return false
}

func (l *Lexer) retreat() {
	if l.position > 0 {
		l.position--
	}

	l.current_char = l.src[l.position]
}

func Tokenise(file string) ([]Token, error) {
	l := new_lexer(file)

	var err error

	for {

		var eof bool
		var token string

		if l.current_char == '\\' {
			l.advance()
			fmt.Println("HELLO")
			if l.current_char == '\\' {
				eof, err = l.scan_comment()
			} else {
				l.retreat()
			}
		}

		if err != nil {
			break
		}

		if eof {
			break
		}

		if l.current_char == ' ' || l.current_char == '\n' || l.current_char == '\t' {
			eof, err = l.scan_ws()
		} else {
			token, eof, err = l.read_word()
			l.tokens = append(l.tokens, Token(token))
		}

		if err != nil {
			break
		}

		if eof {
			break
		}

	}

	return l.tokens, err
}

func (l *Lexer) read_word() (string, bool, error) {
	var buffer []byte

	var err error

	var eof bool

	for i := 0; l.current_char != ' ' && l.current_char != '\n' && l.current_char != '\t'; i++ {
		buffer = append(buffer, l.current_char)
		if l.advance() {
			eof = true
			break
		}
	}

	l.advance()

	return string(buffer), eof, err
}

func (l *Lexer) scan_ws() (bool, error) {
	var eof bool
	var err error

	for l.current_char == ' ' || l.current_char == '\n' || l.current_char == '\t' {
		if l.advance() {
			eof = true
			break
		}
	}

	return eof, err
}

func (l *Lexer) scan_comment() (bool, error) {
	var eof bool
	var err error

	for l.current_char != '\n' {
		if l.advance() {
			eof = true
			break
		}
	}

	return eof, err
}
