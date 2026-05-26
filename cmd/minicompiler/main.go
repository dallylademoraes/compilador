package main

import (
	"bufio"
	"compilador/pkg/compiler"
	"crypto/rand"
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

type cliOptions struct {
	inputFile   string
	interactive bool
	outputPath  string
	phases      phaseFlags
}

func main() {
	cabecalho()

	options, err := parseCLIOptions(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	entrada, origem, err := carregarEntrada(options)
	if err != nil {
		fmt.Printf("Erro ao ler entrada: %v\n", err)
		os.Exit(1)
	}

	executarCompilador(entrada, origem, options.phases, options.outputPath)
}

func cabecalho() {
	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║       MINICOMPILADOR – Disciplina de Compiladores             ║")
	fmt.Println("║       Lexer (Go) + Parser LALR(1) (goyacc)                   ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func parseCLIOptions(args []string) (cliOptions, error) {
	fs := flag.NewFlagSet("minicompiler", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var tokensOnly bool
	var parserOnly bool
	var astOnly bool
	var semanticOnly bool
	var interactive bool
	var outputPath string

	fs.BoolVar(&tokensOnly, "tokens", false, "Mostra apenas a fase léxica")
	fs.BoolVar(&parserOnly, "parser", false, "Mostra apenas a fase sintática")
	fs.BoolVar(&astOnly, "ast", false, "Mostra apenas a AST")
	fs.BoolVar(&semanticOnly, "semantic", false, "Mostra apenas a fase semântica")
	fs.BoolVar(&interactive, "i", false, "Força leitura interativa da entrada via terminal")
	fs.StringVar(&outputPath, "o", "", "Salva o código Python gerado em um arquivo")

	if err := fs.Parse(args); err != nil {
		return cliOptions{}, fmt.Errorf("%v\n\n%s", err, usage())
	}

	restantes := fs.Args()
	if len(restantes) > 1 {
		return cliOptions{}, fmt.Errorf("apenas um arquivo de entrada pode ser informado\n\n%s", usage())
	}
	if interactive && len(restantes) == 1 {
		return cliOptions{}, fmt.Errorf("não use -i junto com arquivo de entrada\n\n%s", usage())
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

	inputFile := ""
	if len(restantes) == 1 {
		inputFile = restantes[0]
	}

	return cliOptions{
		inputFile:   inputFile,
		interactive: interactive,
		outputPath:  outputPath,
		phases:      phases,
	}, nil
}

func parseCLIArgs(args []string) (string, phaseFlags, error) {
	options, err := parseCLIOptions(args)
	if err != nil {
		return "", phaseFlags{}, err
	}

	return options.inputFile, options.phases, nil
}

func usage() string {
	return strings.TrimSpace(`
Uso:
  go run ./cmd/minicompiler [opcoes] [arquivo]

Opções:
  -i            Lê a entrada do terminal (não pode ser usado com arquivo)
  -o arquivo.py Salva o código Python gerado no caminho informado
  --tokens      Mostra apenas a fase léxica
  --parser      Mostra apenas a fase sintática
  --ast         Mostra apenas a AST
  --semantic    Mostra apenas a fase semântica

Sem arquivo informado, o compilador lê a entrada do terminal até Ctrl+D (Linux/macOS) ou Ctrl+Z (Windows).
`)
}

func carregarEntrada(options cliOptions) (string, string, error) {
	if options.interactive || options.inputFile == "" {
		return lerEntrada(), "entrada do usuário", nil
	}

	if options.inputFile != "" {
		conteudo, err := os.ReadFile(options.inputFile)
		if err != nil {
			return "", "", err
		}
		return string(conteudo), fmt.Sprintf("arquivo %s", options.inputFile), nil
	}

	return lerEntrada(), "entrada do usuário", nil
}

func executarCompilador(entrada, origem string, phases phaseFlags, outputPath string) {
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
}

func exportarAST(ast *compiler.Program) {
	outDir := filepath.Join("ast", "imagens")

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Printf("\nAviso: não foi possível criar diretório '%s' (%v).\n", outDir, err)
		return
	}

	pngPath, err := novoCaminhoPNG(outDir)
	if err != nil {
		fmt.Printf("\nAviso: não foi possível gerar nome de arquivo da AST (%v).\n", err)
		return
	}

	if err := compiler.ExportarASTPNG(ast, pngPath); err != nil {
		fmt.Printf("\nAviso: não foi possível gerar PNG da AST (%v).\n", err)
		fmt.Println("Dica: instale o Graphviz e garanta que o comando 'dot' esteja no PATH.")
		return
	}

	dotPath := strings.TrimSuffix(pngPath, filepath.Ext(pngPath)) + ".dot"
	fmt.Printf("\nAST exportada para: %s e %s\n", dotPath, pngPath)
}

func novoCaminhoPNG(dir string) (string, error) {
	for i := 0; i < 10; i++ {
		sufixo, err := sufixoAleatorio(6)
		if err != nil {
			return "", err
		}

		caminho := filepath.Join(dir, fmt.Sprintf("ast-%s.png", sufixo))
		if _, err := os.Stat(caminho); os.IsNotExist(err) {
			return caminho, nil
		} else if err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("não foi possível gerar nome único para AST")
}

func sufixoAleatorio(tamanho int) (string, error) {
	const alfabeto = "abcdefghijklmnopqrstuvwxyz"

	bytes := make([]byte, tamanho)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := range bytes {
		bytes[i] = alfabeto[int(bytes[i])%len(alfabeto)]
	}

	return string(bytes), nil
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
