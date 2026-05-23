package compiler

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

const TOKEN_EOF = 0

// InfoToken guarda o tipo do token e seu valor textual
type InfoToken struct {
	Tipo  int
	Valor string
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

// ProximoToken é chamado pelo parser via interface yyLex
func (l *Lexer) ProximoToken() InfoToken {
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
			return InfoToken{KEYWORD_INT, palavra}
		}
		return InfoToken{ID, palavra}
	}

	// Números
	if unicode.IsDigit(l.atual) {
		var sb strings.Builder
		for !l.fimArq && unicode.IsDigit(l.atual) {
			sb.WriteRune(l.atual)
			l.avancar()
		}
		return InfoToken{NUM, sb.String()}
	}

	// Símbolos
	ch := l.atual
	l.avancar()
	switch ch {
	case '=':
		return InfoToken{ASSIGN, "="}
	case '+':
		return InfoToken{PLUS, "+"}
	case '-':
		return InfoToken{MINUS, "-"}
	case '*':
		return InfoToken{TIMES, "*"}
	case '/':
		return InfoToken{DIV, "/"}
	case ';':
		return InfoToken{SEMI, ";"}
	case '(':
		return InfoToken{LPAREN, "("}
	case ')':
		return InfoToken{RPAREN, ")"}
	}

	return InfoToken{ERROR, string(ch)}
}

// YyLex é a estrutura exigida pelo goyacc
type YyLex struct {
	lexer *Lexer
	lval  *yySymType
	Erros []string
	Program *Program
}

func (y *YyLex) Lex(lval *yySymType) int {
	tok := y.lexer.ProximoToken()
	lval.sval = tok.Valor
	y.lexer.tokens = append(y.lexer.tokens, tok)
	return tok.Tipo
}

func (y *YyLex) Error(s string) {
	y.Erros = append(y.Erros, s)
}

// ExecutarParser coordena a análise sintática
func ExecutarParser(entrada string) (*Program, bool) {
	lexer := NovoLexer(entrada)
	yylex := &YyLex{lexer: lexer}
	resultado := yyParse(yylex)
	if resultado != 0 || len(yylex.Erros) > 0 {
		for _, e := range yylex.Erros {
			fmt.Printf("│   [Erro sintático] %s\n", e)
		}
		return nil, false
	}
	return yylex.Program, true
}

// ImprimirTabelaTokens exibe a tabela léxica formatada
func ImprimirTabelaTokens(entrada string) {
	fmt.Println("┌─────────────────────┬──────────────────────┬────────────────────┐")
	fmt.Printf("│ %-19s │ %-20s │ %-18s │\n", "Lexema", "Token", "Categoria")
	fmt.Println("├─────────────────────┼──────────────────────┼────────────────────┤")

	l := NovoLexer(entrada)
	for {
		tok := l.ProximoToken()
		if tok.Tipo == TOKEN_EOF {
			break
		}
		cat := NomeCategoria(tok.Tipo)
		nomeTok := NomeTipoToken(tok.Tipo)
		fmt.Printf("│ %-19s │ %-20s │ %-18s │\n", tok.Valor, nomeTok, cat)
	}
	fmt.Println("└─────────────────────┴──────────────────────┴────────────────────┘")
}

func NomeCategoria(tipo int) string {
	switch tipo {
	case KEYWORD_INT:
		return "palavra-chave"
	case ID:
		return "identificador"
	case ASSIGN:
		return "atribuição"
	case PLUS, MINUS, TIMES, DIV:
		return "operador"
	case SEMI, LPAREN, RPAREN:
		return "delimitador"
	case NUM:
		return "número"
	default:
		return "erro"
	}
}

func NomeTipoToken(tipo int) string {
	switch tipo {
	case KEYWORD_INT:
		return "<palavra-chave>"
	case ID:
		return "<id>"
	case ASSIGN:
		return "<atrib>"
	case PLUS:
		return "<op,+>"
	case MINUS:
		return "<op,->"
	case TIMES:
		return "<op,*>"
	case DIV:
		return "<op,/>"
	case SEMI:
		return "<delim,;>"
	case LPAREN:
		return "<delim,(>"
	case RPAREN:
		return "<delim,)>"
	case NUM:
		return "<num>"
	default:
		return "<erro>"
	}
}
