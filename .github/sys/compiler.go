// Package sys - Ana GO derleyici modülü
// Lexical analizi, parsing, type checking ve code generation içerir
package sys

import (
	"fmt"
	"strings"
	"sync"
)

// Token - Derleyici token yapısı
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

// TokenType - Token türleri
type TokenType string

const (
	// Literal tokenlar
	TokenEOF    TokenType = "EOF"
	TokenINT    TokenType = "INT"
	TokenSTRING TokenType = "STRING"
	TokenIDENT  TokenType = "IDENT"

	// Operator tokenlar
	TokenPLUS   TokenType = "PLUS"
	TokenMINUS  TokenType = "MINUS"
	TokenMUL    TokenType = "MUL"
	TokenDIV    TokenType = "DIV"
	TokenASSIGN TokenType = "ASSIGN"
	TokenEQ     TokenType = "EQ"
	TokenNE     TokenType = "NE"

	// Anahtar kelimeler
	TokenIF     TokenType = "IF"
	TokenELSE   TokenType = "ELSE"
	TokenFOR    TokenType = "FOR"
	TokenFUNC   TokenType = "FUNC"
	TokenRETURN TokenType = "RETURN"

	// Delimiters
	TokenLPAREN TokenType = "LPAREN"
	TokenRPAREN TokenType = "RPAREN"
	TokenLBRACE TokenType = "LBRACE"
	TokenRBRACE TokenType = "RBRACE"
	TokenSEMI   TokenType = "SEMI"
	TokenCOMMA  TokenType = "COMMA"
)

// Lexer - Lexical analyzer
type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	tokens []Token
	mu     sync.Mutex
}

// NewLexer - Yeni lexer oluştur
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   1,
		col:    1,
		tokens: make([]Token, 0),
	}
}

// Tokenize - İnput'u token'lara ayır
func (l *Lexer) Tokenize() ([]Token, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for l.pos < len(l.input) {
		// Whitespace atla
		if l.isWhitespace(l.current()) {
			if l.current() == '\n' {
				l.line++
				l.col = 1
			} else {
				l.col++
			}
			l.pos++
			continue
		}

		token := l.nextToken()
		if token.Type == TokenEOF {
			break
		}
		l.tokens = append(l.tokens, token)
	}

	l.tokens = append(l.tokens, Token{Type: TokenEOF, Line: l.line, Col: l.col})
	return l.tokens, nil
}

// nextToken - Bir sonraki token'ı oku
func (l *Lexer) nextToken() Token {
	start := l.pos
	startLine := l.line
	startCol := l.col

	ch := l.current()

	// Operatörler ve delimiters
	switch ch {
	case '+':
		l.pos++
		l.col++
		return Token{Type: TokenPLUS, Value: "+", Line: startLine, Col: startCol}
	case '-':
		l.pos++
		l.col++
		return Token{Type: TokenMINUS, Value: "-", Line: startLine, Col: startCol}
	case '*':
		l.pos++
		l.col++
		return Token{Type: TokenMUL, Value: "*", Line: startLine, Col: startCol}
	case '/':
		l.pos++
		l.col++
		return Token{Type: TokenDIV, Value: "/", Line: startLine, Col: startCol}
	case '=':
		l.pos++
		l.col++
		if l.current() == '=' {
			l.pos++
			l.col++
			return Token{Type: TokenEQ, Value: "==", Line: startLine, Col: startCol}
		}
		return Token{Type: TokenASSIGN, Value: "=", Line: startLine, Col: startCol}
	case '!':
		l.pos++
		l.col++
		if l.current() == '=' {
			l.pos++
			l.col++
			return Token{Type: TokenNE, Value: "!=", Line: startLine, Col: startCol}
		}
		return Token{Type: TokenNE, Value: "!", Line: startLine, Col: startCol}
	case '(':
		l.pos++
		l.col++
		return Token{Type: TokenLPAREN, Value: "(", Line: startLine, Col: startCol}
	case ')':
		l.pos++
		l.col++
		return Token{Type: TokenRPAREN, Value: ")", Line: startLine, Col: startCol}
	case '{':
		l.pos++
		l.col++
		return Token{Type: TokenLBRACE, Value: "{", Line: startLine, Col: startCol}
	case '}':
		l.pos++
		l.col++
		return Token{Type: TokenRBRACE, Value: "}", Line: startLine, Col: startCol}
	case ';':
		l.pos++
		l.col++
		return Token{Type: TokenSEMI, Value: ";", Line: startLine, Col: startCol}
	case ',':
		l.pos++
		l.col++
		return Token{Type: TokenCOMMA, Value: ",", Line: startLine, Col: startCol}
	case '"':
		return l.readString(startLine, startCol)
	}

	// Sayılar
	if l.isDigit(ch) {
		return l.readNumber(startLine, startCol)
	}

	// İdentifier ve anahtar kelimeler
	if l.isLetter(ch) {
		return l.readIdent(startLine, startCol)
	}

	l.pos++
	return Token{Type: TokenEOF, Line: startLine, Col: startCol}
}

// readNumber - Sayı oku
func (l *Lexer) readNumber(line, col int) Token {
	start := l.pos
	for l.pos < len(l.input) && l.isDigit(l.current()) {
		l.pos++
		l.col++
	}
	return Token{Type: TokenINT, Value: l.input[start:l.pos], Line: line, Col: col}
}

// readString - String oku
func (l *Lexer) readString(line, col int) Token {
	l.pos++ // "
	l.col++
	start := l.pos

	for l.pos < len(l.input) && l.current() != '"' {
		if l.current() == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.pos++
	}

	value := l.input[start:l.pos]
	if l.pos < len(l.input) {
		l.pos++ // closing "
		l.col++
	}
	return Token{Type: TokenSTRING, Value: value, Line: line, Col: col}
}

// readIdent - İdentifier oku
func (l *Lexer) readIdent(line, col int) Token {
	start := l.pos
	for l.pos < len(l.input) && (l.isLetter(l.current()) || l.isDigit(l.current()) || l.current() == '_') {
		l.pos++
		l.col++
	}

	value := l.input[start:l.pos]
	tokenType := l.getKeywordType(value)
	return Token{Type: tokenType, Value: value, Line: line, Col: col}
}

// getKeywordType - Anahtar kelime türünü belirle
func (l *Lexer) getKeywordType(word string) TokenType {
	switch word {
	case "if":
		return TokenIF
	case "else":
		return TokenELSE
	case "for":
		return TokenFOR
	case "func":
		return TokenFUNC
	case "return":
		return TokenRETURN
	default:
		return TokenIDENT
	}
}

// Helper metodlar
func (l *Lexer) current() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// Parser - Syntax analyzer
type Parser struct {
	tokens  []Token
	current int
	errors  []string
	mu      sync.Mutex
}

// NewParser - Yeni parser oluştur
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errors:  make([]string, 0),
	}
}

// Parse - Token'ları parse et
func (p *Parser) Parse() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for !p.isAtEnd() {
		if err := p.parseStatement(); err != nil {
			p.errors = append(p.errors, err.Error())
			p.advance()
		}
	}

	if len(p.errors) > 0 {
		return fmt.Errorf("parse hatalar:\n%s", strings.Join(p.errors, "\n"))
	}
	return nil
}

// parseStatement - Statement parse et
func (p *Parser) parseStatement() error {
	token := p.peek()

	switch token.Type {
	case TokenIF:
		return p.parseIfStatement()
	case TokenFOR:
		return p.parseForStatement()
	case TokenFUNC:
		return p.parseFuncDeclaration()
	case TokenRETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpression()
	}
}

// parseIfStatement - If statement parse et
func (p *Parser) parseIfStatement() error {
	p.advance() // "if" consume

	// Condition parse et
	if err := p.parseExpression(); err != nil {
		return err
	}

	if p.peek().Type != TokenLBRACE {
		return fmt.Errorf("if statement'da '{' bekleniyor (line %d)", p.peek().Line)
	}

	return p.parseBlock()
}

// parseForStatement - For statement parse et
func (p *Parser) parseForStatement() error {
	p.advance() // "for" consume

	// Initialization parse et
	if p.peek().Type != TokenSEMI {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	p.advance()

	// Condition parse et
	if p.peek().Type != TokenSEMI {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}
	p.advance()

	// Increment parse et
	if p.peek().Type != TokenLBRACE {
		if err := p.parseExpression(); err != nil {
			return err
		}
	}

	return p.parseBlock()
}

// parseFuncDeclaration - Function declaration parse et
func (p *Parser) parseFuncDeclaration() error {
	p.advance() // "func" consume

	// Function adı
	if p.peek().Type != TokenIDENT {
		return fmt.Errorf("function adı bekleniyor (line %d)", p.peek().Line)
	}
	p.advance()

	// Parameters
	if p.peek().Type != TokenLPAREN {
		return fmt.Errorf("'(' bekleniyor (line %d)", p.peek().Line)
	}
	p.advance()

	for p.peek().Type != TokenRPAREN && !p.isAtEnd() {
		p.advance()
	}
	p.advance() // ")" consume

	// Body
	return p.parseBlock()
}

// parseReturnStatement - Return statement parse et
func (p *Parser) parseReturnStatement() error {
	p.advance() // "return" consume
	return p.parseExpression()
}

// parseExpression - Expression parse et
func (p *Parser) parseExpression() error {
	token := p.peek()

	if token.Type == TokenEOF || token.Type == TokenSEMI {
		return nil
	}

	p.advance()
	return nil
}

// parseBlock - Code block parse et
func (p *Parser) parseBlock() error {
	if p.peek().Type != TokenLBRACE {
		return fmt.Errorf("'{' bekleniyor (line %d)", p.peek().Line)
	}
	p.advance()

	depth := 1
	for depth > 0 && !p.isAtEnd() {
		if p.peek().Type == TokenLBRACE {
			depth++
		} else if p.peek().Type == TokenRBRACE {
			depth--
		}
		p.advance()
	}

	return nil
}

// Helper metodlar
func (p *Parser) peek() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.current]
}

func (p *Parser) advance() {
	p.current++
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.peek().Type == TokenEOF
}

// GetErrors - Parse hatalarını döndür
func (p *Parser) GetErrors() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.errors
}

// Compiler - Ana derleyici yapısı
type Compiler struct {
	source   string
	lexer    *Lexer
	parser   *Parser
	optimizer *Optimizer
	errors   []string
	mu       sync.RWMutex
}

// NewCompiler - Yeni derleyici oluştur
func NewCompiler(source string) *Compiler {
	return &Compiler{
		source: source,
		errors: make([]string, 0),
	}
}

// Compile - Kaynağı derle
func (c *Compiler) Compile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Tokenize
	c.lexer = NewLexer(c.source)
	tokens, err := c.lexer.Tokenize()
	if err != nil {
		return fmt.Errorf("tokenize hatası: %w", err)
	}

	// Parse
	c.parser = NewParser(tokens)
	if err := c.parser.Parse(); err != nil {
		c.errors = append(c.errors, err.Error())
		return err
	}

	// Optimize
	c.optimizer = NewOptimizer(c.source)
	if optimized, err := c.optimizer.Optimize(); err != nil {
		c.errors = append(c.errors, fmt.Sprintf("optimize hatası: %v", err))
	} else {
		c.source = optimized
	}

	return nil
}

// GetErrors - Derleyici hatalarını döndür
func (c *Compiler) GetErrors() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.errors
}

// GetOutput - Derlenmiş çıktıyı döndür
func (c *Compiler) GetOutput() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.source
}
