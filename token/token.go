package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL     = "ILLEGAL"
	END_OF_FILE = "EOF"
	IDENTIFIER  = "IDENT"
	INT         = "INT"

	// block comments
	COMMENT_BLOCK_START   = "/*"
	COMMENT_BLOCK_STOP    = "*/"
	COMMENT_DOCUMENTATION = "///"
	COMMENT_LINK          = "//"

	// operators
	OP_ARROW               = "->"
	OP_AT                  = "@"
	OP_BRACKET_ANGLE_LEFT  = "<"
	OP_BRACKET_ANGLE_RIGHT = ">"
	OP_BRACKET_CURLY_LEFT  = "{"
	OP_BRACKET_CURLY_RIGHT = "}"
	OP_BRACKET_PAREN_LEFT  = "("
	OP_BRACKET_PAREN_RIGHT = ")"
	OP_COLON               = ":"
	OP_COMMA               = ","
	OP_EQUAL               = "="
	OP_EXPLICIT_ID         = "%"
	OP_MINUS               = "-"
	OP_PERIOD              = "."
	OP_PLUS                = "+"
	OP_SEMICOLON           = ";"
	OP_SLASH               = "/"
	OP_STAR                = "*"
	OP_UNDERSCORE          = "_"

	// keywords
	KEYWORD_AS          = "AS"
	KEYWORD_BOOL        = "BOOL"
	KEYWORD_BORROW      = "BORROW"
	KEYWORD_CHAR        = "CHAR"
	KEYWORD_CONSTRUCTOR = "CONSTRUCTOR"
	KEYWORD_ENUM        = "ENUM"
	KEYWORD_EXPORT      = "EXPORT"
	KEYWORD_FLAGS       = "FLAGS"
	KEYWORD_FLOAT32     = "FLOAT32"
	KEYWORD_FLOAT64     = "FLOAT64"
	KEYWORD_FROM        = "FROM"
	KEYWORD_FUNC        = "FUNC"
	KEYWORD_FUTURE      = "FUTURE"
	KEYWORD_IMPORT      = "IMPORT"
	KEYWORD_INCLUDE     = "INCLUDE"
	KEYWORD_INTERFACE   = "INTERFACE"
	KEYWORD_LIST        = "LIST"
	KEYWORD_OPTION      = "OPTION"
	KEYWORD_OWN         = "OWN"
	KEYWORD_PACKAGE     = "PACKAGE"
	KEYWORD_RECORD      = "RECORD"
	KEYWORD_RESOURCE    = "RESOURCE"
	KEYWORD_RESULT      = "RESULT"
	KEYWORD_S16         = "S16"
	KEYWORD_S32         = "S32"
	KEYWORD_S64         = "S64"
	KEYWORD_S8          = "S8"
	KEYWORD_STATIC      = "STATIC"
	KEYWORD_STREAM      = "STREAM"
	KEYWORD_STRING      = "STRING"
	KEYWORD_TUPLE       = "TUPLE"
	KEYWORD_TYPE        = "TYPE"
	KEYWORD_U16         = "U16"
	KEYWORD_U32         = "U32"
	KEYWORD_U64         = "U64"
	KEYWORD_U8          = "U8"
	KEYWORD_UNION       = "UNION"
	KEYWORD_USE         = "USE"
	KEYWORD_VARIANT     = "VARIANT"
	KEYWORD_WITH        = "WITH"
	KEYWORD_WORLD       = "WORLD"
)

var keywords = map[string]TokenType{
	"as":          KEYWORD_AS,
	"bool":        KEYWORD_BOOL,
	"borrow":      KEYWORD_BORROW,
	"char":        KEYWORD_CHAR,
	"constructor": KEYWORD_CONSTRUCTOR,
	"enum":        KEYWORD_ENUM,
	"export":      KEYWORD_EXPORT,
	"flags":       KEYWORD_FLAGS,
	"float32":     KEYWORD_FLOAT32,
	"float64":     KEYWORD_FLOAT64,
	"from":        KEYWORD_FROM,
	"func":        KEYWORD_FUNC,
	"future":      KEYWORD_FUTURE,
	"import":      KEYWORD_IMPORT,
	"include":     KEYWORD_INCLUDE,
	"interface":   KEYWORD_INTERFACE,
	"list":        KEYWORD_LIST,
	"option":      KEYWORD_OPTION,
	"own":         KEYWORD_OWN,
	"package":     KEYWORD_PACKAGE,
	"record":      KEYWORD_RECORD,
	"resource":    KEYWORD_RESOURCE,
	"result":      KEYWORD_RESULT,
	"s16":         KEYWORD_S16,
	"s32":         KEYWORD_S32,
	"s64":         KEYWORD_S64,
	"s8":          KEYWORD_S8,
	"static":      KEYWORD_STATIC,
	"stream":      KEYWORD_STREAM,
	"string":      KEYWORD_STRING,
	"tuple":       KEYWORD_TUPLE,
	"type":        KEYWORD_TYPE,
	"u16":         KEYWORD_U16,
	"u32":         KEYWORD_U32,
	"u64":         KEYWORD_U64,
	"u8":          KEYWORD_U8,
	"union":       KEYWORD_UNION,
	"use":         KEYWORD_USE,
	"variant":     KEYWORD_VARIANT,
	"with":        KEYWORD_WITH,
	"world":       KEYWORD_WORLD,
}

func LookupIdentifier(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
