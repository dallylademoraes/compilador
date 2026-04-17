package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Os valores abaixo devem coincidir com os %token declarados no parser.y
const (
	TOKEN_EOF   = 0
	TOKEN_INT   = 57346 // KEYWORD_INT
	TOKEN_ID    = 57347 // ID
	TOKEN_ATRIB = 57348 // ASSIGN
	TOKEN_MAIS  = 57349 // PLUS
	TOKEN_MENOS = 57350 // MINUS
	TOKEN_VEZES = 57351 // TIMES
	TOKEN_PVIRG = 57352 // SEMI
	TOKEN_APAR  = 57353 // LPAREN
	TOKEN_FPAR  = 57354 // RPAREN
	TOKEN_NUM   = 57355 // NUM
	TOKEN_ERRO  = 57356 // ERROR
)

// InfoToken guarda o tipo do token e seu valor textual
type InfoToken struct {
	tipo  int
	valor string
}

// Lexer realiza a análise léxica
type Lexer struct {
	leitor *bufio.Reader
	atual  rune
	fimArq bool
	tokens []InfoToken
}

func NovoLexer(entrada string) *Lexer {
	l := &Lexer{leitor: bufio.NewReader(strings.NewReader(entrada))}
	l.avancar()
	return l
}

func (l *Lexer) avancar() {
	r, _, err := l.leitor.ReadRune()
	if err != nil {
		l.fimArq = true
		l.atual = 0
		return
	}
	l.atual = r
}

// proximoToken é chamado pelo parser via interface yyLex
func (l *Lexer) proximoToken() InfoToken {
	for !l.fimArq && unicode.IsSpace(l.atual) {
		l.avancar()
	}
	if l.fimArq {
		return InfoToken{TOKEN_EOF, "EOF"}
	}

	// Identificadores e palavras-chave
	if unicode.IsLetter(l.atual) || l.atual == '_' {
		var sb strings.Builder
		for !l.fimArq && (unicode.IsLetter(l.atual) || unicode.IsDigit(l.atual) || l.atual == '_') {
			sb.WriteRune(l.atual)
			l.avancar()
		}
		palavra := sb.String()
		if palavra == "int" {
			return InfoToken{TOKEN_INT, palavra}
		}
		return InfoToken{TOKEN_ID, palavra}
	}

	// Números
	if unicode.IsDigit(l.atual) {
		var sb strings.Builder
		for !l.fimArq && unicode.IsDigit(l.atual) {
			sb.WriteRune(l.atual)
			l.avancar()
		}
		return InfoToken{TOKEN_NUM, sb.String()}
	}

	// Símbolos
	ch := l.atual
	l.avancar()
	switch ch {
	case '=':
		return InfoToken{TOKEN_ATRIB, "="}
	case '+':
		return InfoToken{TOKEN_MAIS, "+"}
	case '-':
		return InfoToken{TOKEN_MENOS, "-"}
	case '*':
		return InfoToken{TOKEN_VEZES, "*"}
	case ';':
		return InfoToken{TOKEN_PVIRG, ";"}
	case '(':
		return InfoToken{TOKEN_APAR, "("}
	case ')':
		return InfoToken{TOKEN_FPAR, ")"}
	}

	return InfoToken{TOKEN_ERRO, string(ch)}
}

// Interface yyLexer exigida pelo goyacc

type YyLex struct {
	lexer *Lexer
	lval  *yySymType
	erros []string
}

func (y *YyLex) Lex(lval *yySymType) int {
	tok := y.lexer.proximoToken()
	lval.sval = tok.valor
	y.lexer.tokens = append(y.lexer.tokens, tok)
	return tok.tipo
}

func (y *YyLex) Error(s string) {
	y.erros = append(y.erros, s)
}

// ImprimirTabelaTokens exibe a tabela léxica formatada
func ImprimirTabelaTokens(entrada string) {
	fmt.Println("┌─────────────────────┬──────────────────────┬────────────────────┐")
	fmt.Printf("│ %-19s │ %-20s │ %-18s │\n", "Lexema", "Token", "Categoria")
	fmt.Println("├─────────────────────┼──────────────────────┼────────────────────┤")

	l := NovoLexer(entrada)
	for {
		tok := l.proximoToken()
		if tok.tipo == TOKEN_EOF {
			break
		}
		cat := nomeCategoria(tok.tipo)
		nomeTok := nomeTipoToken(tok.tipo)
		fmt.Printf("│ %-19s │ %-20s │ %-18s │\n", tok.valor, nomeTok, cat)
	}
	fmt.Println("└─────────────────────┴──────────────────────┴────────────────────┘")
}

func nomeCategoria(tipo int) string {
	switch tipo {
	case TOKEN_INT:
		return "palavra-chave"
	case TOKEN_ID:
		return "identificador"
	case TOKEN_ATRIB:
		return "atribuição"
	case TOKEN_MAIS, TOKEN_MENOS, TOKEN_VEZES:
		return "operador"
	case TOKEN_PVIRG, TOKEN_APAR, TOKEN_FPAR:
		return "delimitador"
	case TOKEN_NUM:
		return "número"
	default:
		return "erro"
	}
}

func nomeTipoToken(tipo int) string {
	switch tipo {
	case TOKEN_INT:
		return "<palavra-chave>"
	case TOKEN_ID:
		return "<id>"
	case TOKEN_ATRIB:
		return "<atrib>"
	case TOKEN_MAIS:
		return "<op,+>"
	case TOKEN_MENOS:
		return "<op,->"
	case TOKEN_VEZES:
		return "<op,*>"
	case TOKEN_PVIRG:
		return "<delim,;>"
	case TOKEN_APAR:
		return "<delim,(>"
	case TOKEN_FPAR:
		return "<delim,)>"
	case TOKEN_NUM:
		return "<num>"
	default:
		return "<erro>"
	}
}

func lerEntrada() string {
	fmt.Println("Digite o código (Ctrl+D para finalizar):")
	scanner := bufio.NewScanner(os.Stdin)
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteRune('\n')
	}
	return sb.String()
}