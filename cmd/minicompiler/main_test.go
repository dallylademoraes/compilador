package main

import (
	"compilador/pkg/compiler"
	"fmt"
	"strings"
	"testing"
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
		descricao: "Expressão com divisão",
		codigo:    "int indice = custoTotal / capacidade;",
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

func TestExemplos(t *testing.T) {
	fmt.Println("\n═══════════════════════════════════════════════════════════════")
	fmt.Println("  VALIDAÇÃO DE EXEMPLOS")
	fmt.Println("═══════════════════════════════════════════════════════════════")

	for i, ex := range exemplos {
		fmt.Printf("┌── Exemplo %d: %s\n", i+1, ex.descricao)
		fmt.Printf("│   Código: %s\n│\n", strings.ReplaceAll(ex.codigo, "\n", " | "))

		fmt.Println("│   [Lexer] Tabela de tokens:")
		imprimirTabelaIndentada(ex.codigo)

		fmt.Println("│   [Parser] Derivação:")
		ast, ok := compiler.ExecutarParser(ex.codigo)
		if ok {
			fmt.Println("│   [AST] Estrutura gerada:")
			compiler.ImprimirAST(ast, 4)
		}

		if ok == ex.valido {
			if ex.valido && ast == nil {
				t.Errorf("Exemplo %d: AST não foi gerada para entrada válida", i+1)
			}
			if !ex.valido && ast != nil {
				t.Errorf("Exemplo %d: AST foi gerada para entrada inválida", i+1)
			}

			if ok {
				fmt.Println("│\n└── ✔  ACEITO (esperado: aceito)")
			} else {
				fmt.Println("└── ✔  REJEITADO corretamente (esperado: rejeitar)")
			}
		} else {
			if ok {
				t.Errorf("Exemplo %d: ACEITO mas deveria ter sido rejeitado", i+1)
			} else {
				t.Errorf("Exemplo %d: REJEITADO mas deveria ter sido aceito", i+1)
			}
		}
		fmt.Println()
	}
}

func TestAnaliseSemantica(t *testing.T) {
	t.Parallel()

	testes := []struct {
		nome              string
		codigo            string
		esperaSucesso     bool
		errosEsperados    []string
		simbolosEsperados []string
	}{
		{
			nome:              "programa semanticamente válido",
			codigo:            "int x = 10;\nint y = x + 1;\ny = y + x;",
			esperaSucesso:     true,
			simbolosEsperados: []string{"x", "y"},
		},
		{
			nome:           "uso antes da declaração",
			codigo:         "int y = x + 1;",
			esperaSucesso:  false,
			errosEsperados: []string{"[Erro semântico] Variável 'x' usada antes de ser declarada"},
		},
		{
			nome:           "atribuição sem declaração",
			codigo:         "x = 42;",
			esperaSucesso:  false,
			errosEsperados: []string{"[Erro semântico] Variável 'x' usada antes de ser declarada"},
		},
		{
			nome:           "redeclaração de variável",
			codigo:         "int x = 1;\nint x = 2;",
			esperaSucesso:  false,
			errosEsperados: []string{"[Erro semântico] Variável 'x' já foi declarada"},
		},
	}

	for _, teste := range testes {
		teste := teste
		t.Run(teste.nome, func(t *testing.T) {
			ast, ok := compiler.ExecutarParser(teste.codigo)
			if !ok {
				t.Fatalf("parser rejeitou um caso de teste semântico válido para sintaxe: %q", teste.codigo)
			}

			tabela, erros := compiler.AnalisarSemantica(ast)
			if (len(erros) == 0) != teste.esperaSucesso {
				t.Fatalf("resultado semântico inesperado: erros=%v", erros)
			}

			for _, erroEsperado := range teste.errosEsperados {
				if !contains(erros, erroEsperado) {
					t.Fatalf("erro esperado não encontrado: %q em %v", erroEsperado, erros)
				}
			}

			for _, nome := range teste.simbolosEsperados {
				if _, ok := tabela.Buscar(nome); !ok {
					t.Fatalf("símbolo %q não foi registrado na tabela", nome)
				}
			}
		})
	}
}

func TestParseCLIArgs(t *testing.T) {
	t.Parallel()

	testes := []struct {
		nome          string
		args          []string
		esperaArquivo string
		esperaErro    bool
		fases         phaseFlags
	}{
		{
			nome:  "sem argumentos mostra todas as fases",
			args:  nil,
			fases: phaseFlags{tokens: true, parser: true, ast: true, semantic: true},
		},
		{
			nome:          "arquivo posicional",
			args:          []string{"programa.txt"},
			esperaArquivo: "programa.txt",
			fases:         phaseFlags{tokens: true, parser: true, ast: true, semantic: true},
		},
		{
			nome:          "flags de fase com arquivo",
			args:          []string{"--tokens", "--semantic", "programa.txt"},
			esperaArquivo: "programa.txt",
			fases:         phaseFlags{tokens: true, semantic: true},
		},
		{
			nome:       "mais de um arquivo",
			args:       []string{"a.txt", "b.txt"},
			esperaErro: true,
		},
	}

	for _, teste := range testes {
		teste := teste
		t.Run(teste.nome, func(t *testing.T) {
			arquivo, fases, err := parseCLIArgs(teste.args)
			if teste.esperaErro {
				if err == nil {
					t.Fatal("era esperado erro ao interpretar os argumentos")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado ao interpretar os argumentos: %v", err)
			}

			if arquivo != teste.esperaArquivo {
				t.Fatalf("arquivo inesperado: esperado %q, obtido %q", teste.esperaArquivo, arquivo)
			}

			if fases != teste.fases {
				t.Fatalf("seleção de fases inesperada: esperado %+v, obtido %+v", teste.fases, fases)
			}
		})
	}
}

func imprimirTabelaIndentada(entrada string) {
	l := compiler.NovoLexer(entrada)
	fmt.Println("│   ┌───────────────┬─────────────────┬────────────────┐")
	fmt.Printf("│   │ %-13s │ %-15s │ %-14s │\n", "Lexema", "Token", "Categoria")
	fmt.Println("│   ├───────────────┼─────────────────┼────────────────┤")
	for {
		tok := l.ProximoToken()
		if tok.Tipo == compiler.TOKEN_EOF {
			break
		}
		fmt.Printf("│   │ %-13s │ %-15s │ %-14s │\n",
			tok.Valor, compiler.NomeTipoToken(tok.Tipo), compiler.NomeCategoria(tok.Tipo))
	}
	fmt.Println("│   └───────────────┴─────────────────┴────────────────┘")
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}
