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

// TODO: missing floats
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '@':
		tok = token.Token{Type: token.OP_AT, Literal: string(l.ch)}
	case '<':
		tok = token.Token{Type: token.OP_BRACKET_ANGLE_LEFT, Literal: string(l.ch)}
	case '>':
		tok = token.Token{Type: token.OP_BRACKET_ANGLE_RIGHT, Literal: string(l.ch)}
	case '{':
		tok = token.Token{Type: token.OP_BRACKET_CURLY_LEFT, Literal: string(l.ch)}
	case '}':
		tok = token.Token{Type: token.OP_BRACKET_CURLY_RIGHT, Literal: string(l.ch)}
	case '(':
		tok = token.Token{Type: token.OP_BRACKET_PAREN_LEFT, Literal: string(l.ch)}
	case ')':
		tok = token.Token{Type: token.OP_BRACKET_PAREN_RIGHT, Literal: string(l.ch)}
	case ':':
		tok = token.Token{Type: token.OP_COLON, Literal: string(l.ch)}
	case ',':
		tok = token.Token{Type: token.OP_COMMA, Literal: string(l.ch)}
	case '=':
		tok = token.Token{Type: token.OP_EQUAL, Literal: string(l.ch)}
	case '%':
		tok = token.Token{Type: token.OP_EXPLICIT_ID, Literal: string(l.ch)}
	case '-':
		p := l.peek()
		if p == '>' {
			tok = token.Token{Type: token.OP_ARROW, Literal: string(l.ch) + string(p)}
			l.readChar()
		} else {
			tok = token.Token{Type: token.OP_MINUS, Literal: string(l.ch)}
		}
	case '.':
		tok = token.Token{Type: token.OP_PERIOD, Literal: string(l.ch)}
	case '+':
		tok = token.Token{Type: token.OP_PLUS, Literal: string(l.ch)}
	case ';':
		tok = token.Token{Type: token.OP_SEMICOLON, Literal: string(l.ch)}
	case '*':
		tok = token.Token{Type: token.OP_STAR, Literal: string(l.ch)}
	case '_':
		tok = token.Token{Type: token.OP_UNDERSCORE, Literal: string(l.ch)}
	case 0:
		tok = token.Token{Type: token.END_OF_FILE, Literal: string("")}
	default:
		if unicode.IsLetter(rune(l.ch)) {
			lit := l.readIdentifier()
			switch lit {
			case "u":
				num := l.peekNum()

				switch num {
				case 8:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_U8, Literal: string(lit) + n}
				case 16:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_U16, Literal: string(lit) + n}
				case 32:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_U32, Literal: string(lit) + n}
				case 64:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_U64, Literal: string(lit) + n}
				default:
					return token.Token{Type: token.LookupIdentifier(lit), Literal: lit}
				}
			case "s":
				num := l.peekNum()
				switch num {
				case 8:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_S8, Literal: string(lit) + n}
				case 16:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_S16, Literal: string(lit) + n}
				case 32:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_S32, Literal: string(lit) + n}
				case 64:
					n := l.readNumber()
					return token.Token{Type: token.KEYWORD_S64, Literal: string(lit) + n}
				default:
					return token.Token{Type: token.LookupIdentifier(lit), Literal: lit}
				}
			default:
				return token.Token{Type: token.LookupIdentifier(lit), Literal: lit}
			}
		} else if unicode.IsDigit(rune(l.ch)) {
			return token.Token{Type: token.INT, Literal: l.readNumber()}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
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

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
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

func (l *Lexer) readIdentifier() string {
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
