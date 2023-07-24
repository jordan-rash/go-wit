package lexer

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/jordan-rash/go-wit/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := new(Lexer)
	l.input = input
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()

	switch l.ch {
	case '@':
		return token.Token{Type: token.OP_AT, Literal: string(l.readChar())}
	case '<':
		return token.Token{Type: token.OP_BRACKET_ANGLE_LEFT, Literal: string(l.readChar())}
	case '>':
		return token.Token{Type: token.OP_BRACKET_ANGLE_RIGHT, Literal: string(l.readChar())}
	case '{':
		return token.Token{Type: token.OP_BRACKET_CURLY_LEFT, Literal: string(l.readChar())}
	case '}':
		return token.Token{Type: token.OP_BRACKET_CURLY_RIGHT, Literal: string(l.readChar())}
	case '(':
		return token.Token{Type: token.OP_BRACKET_PAREN_LEFT, Literal: string(l.readChar())}
	case ')':
		return token.Token{Type: token.OP_BRACKET_PAREN_RIGHT, Literal: string(l.readChar())}
	case ':':
		return token.Token{Type: token.OP_COLON, Literal: string(l.readChar())}
	case ',':
		return token.Token{Type: token.OP_COMMA, Literal: string(l.readChar())}
	case '=':
		return token.Token{Type: token.OP_EQUAL, Literal: string(l.readChar())}
	case '%':
		return token.Token{Type: token.OP_EXPLICIT_ID, Literal: string(l.readChar())}
	case '-':
		p := l.peek()
		if p == '>' {
			return token.Token{Type: token.OP_ARROW, Literal: string(l.readChar()) + string(l.readChar())}
		} else {
			return token.Token{Type: token.OP_MINUS, Literal: string(l.readChar())}
		}
	case '.':
		return token.Token{Type: token.OP_PERIOD, Literal: string(l.readChar())}
	case '+':
		return token.Token{Type: token.OP_PLUS, Literal: string(l.readChar())}
	case ';':
		return token.Token{Type: token.OP_SEMICOLON, Literal: string(l.readChar())}
	case '*':
		x := l.readChar()
		return token.Token{Type: token.OP_STAR, Literal: string(x)}
	case '_':
		x := l.readChar()
		return token.Token{Type: token.OP_UNDERSCORE, Literal: string(x)}
	case 0:
		return token.Token{Type: token.END_OF_FILE, Literal: string("EOF")}
	case 'u':
		switch l.peekNum() {
		case 8:
			return token.Token{Type: token.KEYWORD_U8, Literal: "u" + l.readNumber()}
		case 16:
			return token.Token{Type: token.KEYWORD_U16, Literal: "u" + l.readNumber()}
		case 32:
			return token.Token{Type: token.KEYWORD_U32, Literal: "u" + l.readNumber()}
		case 64:
			return token.Token{Type: token.KEYWORD_U64, Literal: "u" + l.readNumber()}
		}
	case 's':
		switch l.peekNum() {
		case 8:
			return token.Token{Type: token.KEYWORD_S8, Literal: "s" + l.readNumber()}
		case 16:
			return token.Token{Type: token.KEYWORD_S16, Literal: "s" + l.readNumber()}
		case 32:
			return token.Token{Type: token.KEYWORD_S32, Literal: "s" + l.readNumber()}
		case 64:
			return token.Token{Type: token.KEYWORD_S64, Literal: "s" + l.readNumber()}
		}
	case 'f':
		if l.peekString() == "float" {
			l.readString()
			switch l.peekNum() {
			case 32:
				return token.Token{Type: token.KEYWORD_FLOAT32, Literal: "float" + l.readNumber()}
			case 64:
				return token.Token{Type: token.KEYWORD_FLOAT64, Literal: "float" + l.readNumber()}
			}
		}
		// case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		// 	if l.peek() == 0 || l.peek() != '-' || unicode.IsSpace(rune(l.ch)) {
		// 		n := l.readNumber()
		// 		return token.Token{Type: token.INT, Literal: n}
		// 	}
	}

	lit := l.readIdentifier()
	i, err := strconv.ParseInt(lit, 10, 0)
	if err == nil {
		return token.Token{Type: token.INT, Literal: strconv.FormatInt(i, 10)}
	}

	return token.Token{Type: token.LookupIdentifier(lit), Literal: lit}
}

func (l *Lexer) peekNum() int {
	origPos := l.position

	sb := strings.Builder{}
	sb.WriteByte(l.ch)

	for unicode.IsDigit(rune(l.peek())) {
		b := l.peek()
		sb.WriteByte(b)
		l.readPosition++
	}

	i, _ := strconv.ParseInt(sb.String(), 10, 0)

	l.position = origPos
	return int(i)
}

func (l *Lexer) peekString() string {
	origPos := l.position

	sb := strings.Builder{}
	sb.WriteByte(l.ch)

	for unicode.IsLetter(rune(l.peek())) {
		b := l.peek()
		sb.WriteByte(b)
		l.readPosition++
	}

	l.position = origPos
	return sb.String()
}

func (l *Lexer) readChar() byte {
	var ret byte
	if l.readPosition >= len(l.input) {
		l.ch = 0
		ret = 0
	} else {
		ret = l.ch
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++

	return ret
}

func (l *Lexer) readNumber() string {
	pos := l.position
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readString() string {
	pos := l.position
	for unicode.IsLetter(rune(l.ch)) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

// kabab case
func (l *Lexer) peekIdentifier() string {
	origPos := l.position

	sb := strings.Builder{}
	sb.WriteByte(l.ch)

	for l.peek() == '-' || unicode.IsDigit(rune(l.peek())) || unicode.IsLetter(rune(l.peek())) {
		b := l.peek()
		sb.WriteByte(b)
		l.readPosition++
	}

	l.position = origPos
	return sb.String()
}

func (l *Lexer) readIdentifier() string {
	pos := l.position

	for l.peek() == '-' || unicode.IsDigit(rune(l.peek())) || unicode.IsLetter(rune(l.peek())) {
		l.readChar()
	}

	l.readChar()
	return l.input[pos:l.position]
}
