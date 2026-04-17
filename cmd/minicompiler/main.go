package main

import (
	"bufio"
	"compilador/pkg/compiler"
	"fmt"
	"os"
	"strings"
)

func main() {
	cabecalho()

	if len(os.Args) > 1 && os.Args[1] == "-i" {
		entrada := lerEntrada()
		executarCompilador(entrada, "entrada do usuário")
		return
	}

	fmt.Println("Use '-i' para entrar no modo interativo ou 'go test ./...' para rodar a validação de exemplos.")
}

func cabecalho() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║       MINICOMPILADOR – Disciplina de Compiladores             ║")
	fmt.Println("║       Lexer (Go) + Parser LALR(1) (goyacc)                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func executarCompilador(entrada, origem string) {
	fmt.Printf("\n── Código de %s ──\n%s\n\n", origem, entrada)
	fmt.Println("── Fase 1: Análise Léxica ──")
	compiler.ImprimirTabelaTokens(entrada)
	fmt.Println("\n── Fase 2: Análise Sintática ──")
	ok := compiler.ExecutarParser(entrada)
	if !ok {
		fmt.Println("\n✗  Programa rejeitado pela gramática.")
	}
}

func lerEntrada() string {
	fmt.Println("Digite o código (Pressione Enter e depois Ctrl+D para finalizar no Linux ou Ctrl+Z no Windows):")
	scanner := bufio.NewScanner(os.Stdin)
	var sb strings.Builder
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteRune('\n')
	}
	return sb.String()
}
