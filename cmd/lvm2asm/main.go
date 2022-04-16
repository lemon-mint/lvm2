package main

import "fmt"

type TokenType string

const (
	TokenType_EOF        = "EOF"
	TokenType_NUMBER     = "NUMBER"
	TokenType_STRING     = "STRING"
	TokenType_IDENTIFIER = "IDENTIFIER"
	TokenType_REGISTER   = "REGISTER"
	TokenType_ADD        = "ADD"
	TokenType_SUB        = "SUB"
	TokenType_COMMA      = "COMMA"
	TokenType_COMMENT    = "COMMENT"

	TokenType_LBRACKET = "LBRACKET"
	TokenType_RBRACKET = "RBRACKET"

	TokenType_INSTRUCTION = "INSTRUCTION"
	TokenType_COMPILER    = "COMPILER"
)

type Token struct {
	Type  TokenType
	Value string

	Line int
	Col  int

	Pos int
}

func NewToken(type_ TokenType, value_ string, line_ int, col_ int, pos_ int) *Token {
	return &Token{
		Type:  type_,
		Value: value_,

		Line: line_,
		Col:  col_,

		Pos: pos_,
	}
}

type Lexer struct {
	data []rune
	pos  int

	line int
	col  int

	state map[string]interface{}

	char      rune
	isEscaped bool

	lastToken *Token

	tokens []*Token

	err error
}

func (l *Lexer) ReadChar() bool {
	if l.pos >= len(l.data) {
		l.char = 0
		return false
	}

	l.char = l.data[l.pos]
	l.pos++

	if l.char == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}

	return true
}

func (l *Lexer) SkipWhitespace() bool {
	for l.char != 0 && (l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r') {
		if !l.ReadChar() {
			return false
		}
	}
	return true
}

var asm_keywords = map[string]bool{
	"NOP":     true,
	"ADD":     true,
	"SUB":     true,
	"MUL":     true,
	"DIV":     true,
	"MOD":     true,
	"AND":     true,
	"OR":      true,
	"XOR":     true,
	"NOT":     true,
	"SHL":     true,
	"SHR":     true,
	"CMP":     true,
	"JMP":     true,
	"JG":      true,
	"JL":      true,
	"JE":      true,
	"JNE":     true,
	"JGE":     true,
	"JLE":     true,
	"LOAD":    true,
	"LOADH":   true,
	"LOADB":   true,
	"STORE":   true,
	"STOREH":  true,
	"STOREB":  true,
	"MOV":     true,
	"MOVH":    true,
	"MOVB":    true,
	"PUSH":    true,
	"POP":     true,
	"CALL":    true,
	"RET":     true,
	"SYSCALL": true,
}

var compiler_keywords = map[string]bool{
	"SECTION": true,
	"ALIGN":   true,
	"CHARS":   true,
	"LABEL":   true,
}

var register_keywords = map[string]bool{
	"$R0":    true,
	"$R1":    true,
	"$R2":    true,
	"$R3":    true,
	"$R4":    true,
	"$R5":    true,
	"$R6":    true,
	"$R7":    true,
	"$R8":    true,
	"$R9":    true,
	"$R10":   true,
	"$R11":   true,
	"$R12":   true,
	"$R13":   true,
	"$R14":   true,
	"$R15":   true,
	"$R16":   true,
	"$R17":   true,
	"$R18":   true,
	"$R19":   true,
	"$R20":   true,
	"$R21":   true,
	"$R22":   true,
	"$R23":   true,
	"$R24":   true,
	"$R25":   true,
	"$R26":   true,
	"$R27":   true,
	"$R28":   true,
	"$R29":   true,
	"$R30":   true,
	"$R31":   true,
	"$SYS32": true,
	"$SYS33": true,
	"$SYS34": true,
	"$SYS35": true,
	"$SYS36": true,
	"$SYS37": true,
	"$SYS38": true,
	"$SYS39": true,
	"$SYS40": true,
	"$SYS41": true,
	"$SYS42": true,
	"$SYS43": true,
	"$SYS44": true,
	"$SYS45": true,
	"$SYS46": true,
	"$SYS47": true,
	"$SYS48": true,
	"$SYS49": true,
	"$SYS50": true,
	"$SYS51": true,
	"$SYS52": true,
	"$SYS53": true,
	"$SYS54": true,
	"$SYS55": true,
	"$SYS56": true,
	"$SYS57": true,
	"$SYS58": true,
	"$SYS59": true,
	"$SYS60": true,
	"$SYS61": true,
	"$SYS62": true,
	"$SYS63": true,
	"$PC":    true,
	"$SP":    true,
	"$SB":    true,
}

func LookupKeyword(keyword string) TokenType {
	if asm_keywords[keyword] {
		return TokenType_INSTRUCTION
	} else if compiler_keywords[keyword] {
		return TokenType_COMPILER
	} else if register_keywords[keyword] {
		return TokenType_REGISTER
	}
	return TokenType_IDENTIFIER
}

func NewLexer(data []rune) (*Lexer, error) {
	sl := &Lexer{
		data: data,
		pos:  0,

		line: 1,
		col:  0,

		state: make(map[string]interface{}),
	}

	if !sl.ReadChar() {
		return nil, fmt.Errorf("empty input")
	}
	return sl, nil
}

func main() {

}
