package main

import (
	"bufio"
	"compilador/pkg/compiler"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cabecalho()

	interactive := false
	outputPath := ""

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-i" {
			interactive = true
		} else if os.Args[i] == "-o" && i+1 < len(os.Args) {
			outputPath = os.Args[i+1]
			i++
		}
	}

	if interactive {
		entrada := lerEntrada()
		executarCompilador(entrada, "entrada do usuário", outputPath)
		return
	}

	fmt.Println("Use '-i' para entrar no modo interativo ou 'go test ./...' para rodar a validação de exemplos.")
	fmt.Println("Opcional: use '-o saida.py' para salvar o código Python gerado.")
}

func cabecalho() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║       MINICOMPILADOR – Disciplina de Compiladores             ║")
	fmt.Println("║       Lexer (Go) + Parser LALR(1) (goyacc)                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func executarCompilador(entrada, origem, outputPath string) {
	fmt.Printf("\n── Código de %s ──\n%s\n\n", origem, entrada)
	
	fmt.Println("── Fase 1: Análise Léxica ──")
	compiler.ImprimirTabelaTokens(entrada)
	
	fmt.Println("\n── Fase 2: Análise Sintática ──")
	ast, ok := compiler.ExecutarParser(entrada)
	if !ok {
		fmt.Println("\n✗  Programa rejeitado pela gramática.")
		return
	}

	fmt.Println("\n── Fase 3: AST ──")
	compiler.ImprimirAST(ast, 0)

	// Fase 5: Código Intermediário
	fmt.Println("\n── Fase 5: Código Intermediário (3AC) ──")
	gerador := &compiler.Gerador3AC{}
	for _, stmt := range ast.Statements {
		gerador.VisitNode(stmt)
	}
	for _, inst := range gerador.Instruction_list {
		if inst.Operador == "=" {
			fmt.Printf("%s = %s\n", inst.Result_addr, inst.Var_um)
		} else {
			fmt.Printf("%s = %s %s %s\n", inst.Result_addr, inst.Var_um, inst.Operador, inst.Var_dois)
		}
	}

	// Fase 6: Geração de Código Python
	fmt.Println("\n── Fase 6: Geração de Código Python ──")
	codigoPython := compiler.GerarPython(gerador.Instruction_list)
	fmt.Println(codigoPython)

	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(codigoPython), 0644); err != nil {
			fmt.Printf("\nErro ao salvar arquivo de saída: %v\n", err)
		} else {
			fmt.Printf("\nCódigo Python salvo em: %s\n", outputPath)
		}
	}

	outDir := "ast"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Printf("\nAviso: não foi possível criar diretório '%s' (%v).\n", outDir, err)
		return
	}

	pngPath := filepath.Join(outDir, "ast.png")
	if err := compiler.ExportarASTPNG(ast, pngPath); err != nil {
		fmt.Printf("\nAviso: não foi possível gerar PNG da AST (%v).\n", err)
		fmt.Println("Dica: instale o Graphviz e garanta que o comando 'dot' esteja no PATH.")
	} else {
		dotPath := filepath.Join(outDir, "ast.dot")
		fmt.Printf("\nAST exportada para: %s e %s\n", dotPath, pngPath)
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
