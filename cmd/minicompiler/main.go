package main

import (
	"bufio"
	"compilador/pkg/compiler"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type phaseFlags struct {
	tokens   bool
	parser   bool
	ast      bool
	semantic bool
}

func main() {
	cabecalho()

	inputFile, phases, err := parseCLIArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	entrada, origem, err := carregarEntrada(inputFile)
	if err != nil {
		fmt.Printf("Erro ao ler entrada: %v\n", err)
		os.Exit(1)
	}

	executarCompilador(entrada, origem, phases)
}

func cabecalho() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║       MINICOMPILADOR – Disciplina de Compiladores             ║")
	fmt.Println("║       Lexer (Go) + Parser LALR(1) (goyacc)                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func parseCLIArgs(args []string) (string, phaseFlags, error) {
	fs := flag.NewFlagSet("minicompiler", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var tokensOnly bool
	var parserOnly bool
	var astOnly bool
	var semanticOnly bool

	fs.BoolVar(&tokensOnly, "tokens", false, "Mostra apenas a fase léxica")
	fs.BoolVar(&parserOnly, "parser", false, "Mostra apenas a fase sintática")
	fs.BoolVar(&astOnly, "ast", false, "Mostra apenas a AST")
	fs.BoolVar(&semanticOnly, "semantic", false, "Mostra apenas a fase semântica")

	if err := fs.Parse(args); err != nil {
		return "", phaseFlags{}, fmt.Errorf("%v\n\n%s", err, usage())
	}

	restantes := fs.Args()
	if len(restantes) > 1 {
		return "", phaseFlags{}, fmt.Errorf("apenas um arquivo de entrada pode ser informado\n\n%s", usage())
	}

	phases := phaseFlags{
		tokens:   tokensOnly,
		parser:   parserOnly,
		ast:      astOnly,
		semantic: semanticOnly,
	}
	if !tokensOnly && !parserOnly && !astOnly && !semanticOnly {
		phases = phaseFlags{tokens: true, parser: true, ast: true, semantic: true}
	}

	if len(restantes) == 1 {
		return restantes[0], phases, nil
	}
	return "", phases, nil
}

func usage() string {
	return strings.TrimSpace(`
Uso:
  go run ./cmd/minicompiler
  go run ./cmd/minicompiler nomeDoArquivo.txt
  go run ./cmd/minicompiler [--tokens] [--parser] [--ast] [--semantic] [arquivo]

Sem arquivo informado, o compilador lê a entrada do terminal até Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows).
`)
}

func carregarEntrada(inputFile string) (string, string, error) {
	if inputFile != "" {
		conteudo, err := os.ReadFile(inputFile)
		if err != nil {
			return "", "", err
		}
		return string(conteudo), fmt.Sprintf("arquivo %s", inputFile), nil
	}

	return lerEntrada(), "entrada do usuário", nil
}

func executarCompilador(entrada, origem string, phases phaseFlags) {
	fmt.Printf("\n── Código de %s ──\n%s\n\n", origem, entrada)

	if phases.tokens {
		fmt.Println("── Fase 1: Análise Léxica ──")
		compiler.ImprimirTabelaTokens(entrada)
	}
	if !phases.parser && !phases.ast && !phases.semantic {
		return
	}

	ast, errosSintaticos := compiler.ParsearPrograma(entrada)
	if phases.parser {
		fmt.Println("\n── Fase 2: Análise Sintática ──")
	}
	if len(errosSintaticos) > 0 {
		if !phases.parser {
			fmt.Println("── Fase 2: Análise Sintática ──")
		}
		for _, erro := range errosSintaticos {
			fmt.Printf("│   [Erro sintático] %s\n", erro)
		}
		fmt.Println("\n✗  Programa rejeitado pela gramática.")
		return
	}
	if phases.parser {
		fmt.Println("✔  Programa aceito pela gramática.")
	}
	if !phases.ast && !phases.semantic {
		return
	}

	if phases.ast {
		fmt.Println("\n── Fase 3: AST ──")
		compiler.ImprimirAST(ast, 0)
		exportarAST(ast)
	}
	if !phases.semantic {
		return
	}

	tabela, errosSemanticos := compiler.AnalisarSemantica(ast)
	if phases.semantic {
		fmt.Println("\n── Fase 4: Análise Semântica ──")
	}
	if len(errosSemanticos) > 0 {
		if !phases.semantic {
			fmt.Println("── Fase 4: Análise Semântica ──")
		}
		for _, erro := range errosSemanticos {
			fmt.Println(erro)
		}
		fmt.Println("\n✗  Programa rejeitado semanticamente.")
		return
	}
	if phases.semantic {
		fmt.Println("✔  Análise semântica concluída sem erros.")
		fmt.Println("\nTabela de símbolos:")
		tabela.Imprimir()
	}
}

func exportarAST(ast *compiler.Program) {
	outDir := "ast"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Printf("\nAviso: não foi possível criar diretório '%s' (%v).\n", outDir, err)
		return
	}

	pngPath := filepath.Join(outDir, "ast.png")
	if err := compiler.ExportarASTPNG(ast, pngPath); err != nil {
		fmt.Printf("\nAviso: não foi possível gerar PNG da AST (%v).\n", err)
		fmt.Println("Dica: instale o Graphviz e garanta que o comando 'dot' esteja no PATH.")
		return
	}

	dotPath := filepath.Join(outDir, "ast.dot")
	fmt.Printf("\nAST exportada para: %s e %s\n", dotPath, pngPath)
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
