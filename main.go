package main

import (
	"fmt"
	"os"
	"strings"
)

var exemplos = []struct {
	descricao string
	codigo    string
	valido    bool
}{
	{
		descricao: "Declaração com atribuição simples",
		codigo:    "int soma = a + b;",
		valido:    true,
	},
	{
		descricao: "Declaração com expressão composta",
		codigo:    "int resultado = a + b + c;",
		valido:    true,
	},
	{
		descricao: "Atribuição sem declaração",
		codigo:    "x = 42;",
		valido:    true,
	},
	{
		descricao: "Expressão com multiplicação",
		codigo:    "int total = a + b * c;",
		valido:    true,
	},
	{
		descricao: "Expressão com parênteses",
		codigo:    "int val = (a + b) * c;",
		valido:    true,
	},
	{
		descricao: "Múltiplas declarações",
		codigo:    "int x = a + b;\nint y = x + 1;",
		valido:    true,
	},
	{
		descricao: "ERRO: ponto-e-vírgula faltando",
		codigo:    "int soma = a + b",
		valido:    false,
	},
	{
		descricao: "ERRO: operador inválido",
		codigo:    "int x = a @ b;",
		valido:    false,
	},
}

func main() {
	cabecalho()

	if len(os.Args) > 1 && os.Args[1] == "-i" {
		entrada := lerEntrada()
		executarCompilador(entrada, "entrada do usuário")
		return
	}

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("  VALIDAÇÃO DE EXEMPLOS")
	fmt.Println("═══════════════════════════════════════════════════════════════\n")

	aprovados, reprovados := 0, 0

	for i, ex := range exemplos {
		fmt.Printf("┌── Exemplo %d: %s\n", i+1, ex.descricao)
		fmt.Printf("│   Código: %s\n│\n", strings.ReplaceAll(ex.codigo, "\n", " | "))

		fmt.Println("│   [Lexer] Tabela de tokens:")
		imprimirTabelaIndentada(ex.codigo)

		fmt.Println("│   [Parser] Derivação:")
		ok := executarParser(ex.codigo)

		if ok == ex.valido {
			if ok {
				fmt.Println("│")
				fmt.Println("└── ✔  ACEITO (esperado: aceito)")
			} else {
				fmt.Println("└── ✔  REJEITADO corretamente (esperado: rejeitar)")
			}
			aprovados++
		} else {
			if ok {
				fmt.Println("└── ✗  ACEITO mas deveria ter sido rejeitado")
			} else {
				fmt.Println("└── ✗  REJEITADO mas deveria ter sido aceito")
			}
			reprovados++
		}
		fmt.Println()
	}

	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
	fmt.Printf("  Resultado: %d/%d exemplos corretos\n", aprovados, len(exemplos))
	fmt.Printf("═══════════════════════════════════════════════════════════════\n")
}

func cabecalho() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║       MINICOMPILADOR – Disciplina de Compiladores             ║")
	fmt.Println("║       Lexer (Go) + Parser LALR(1) (goyacc)                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func executarParser(entrada string) bool {
	lexer := NovoLexer(entrada)
	yylex := &YyLex{lexer: lexer}
	resultado := yyParse(yylex)
	if resultado != 0 || len(yylex.erros) > 0 {
		for _, e := range yylex.erros {
			fmt.Printf("│   [Erro sintático] %s\n", e)
		}
		return false
	}
	return true
}

func executarCompilador(entrada, origem string) {
	fmt.Printf("\n── Código de %s ──\n%s\n\n", origem, entrada)
	fmt.Println("── Fase 1: Análise Léxica ──")
	ImprimirTabelaTokens(entrada)
	fmt.Println("\n── Fase 2: Análise Sintática ──")
	ok := executarParser(entrada)
	if !ok {
		fmt.Println("\n✗  Programa rejeitado pela gramática.")
	}
}

func imprimirTabelaIndentada(entrada string) {
	l := NovoLexer(entrada)
	fmt.Println("│   ┌───────────────┬─────────────────┬────────────────┐")
	fmt.Printf("│   │ %-13s │ %-15s │ %-14s │\n", "Lexema", "Token", "Categoria")
	fmt.Println("│   ├───────────────┼─────────────────┼────────────────┤")
	for {
		tok := l.proximoToken()
		if tok.tipo == TOKEN_EOF {
			break
		}
		fmt.Printf("│   │ %-13s │ %-15s │ %-14s │\n",
			tok.valor, nomeTipoToken(tok.tipo), nomeCategoria(tok.tipo))
	}
	fmt.Println("│   └───────────────┴─────────────────┴────────────────┘")
}