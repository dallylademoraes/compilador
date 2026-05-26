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
	valido    bool}{
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

			// --- INTEGRAÇÃO DO CÓDIGO INTERMEDIÁRIO (3AC) ---
			fmt.Println("│\n│   [3AC] Código Intermediário gerado:")
			gerador := &compiler.Gerador3AC{}

			// Varre as linhas da AST
			for _, stmt := range ast.Statements {
				gerador.VisitNode(stmt)
			}

			// Imprime as instruções formatadas
			for _, inst := range gerador.Instruction_list {
				if inst.Operador == "=" {
					fmt.Printf("│       %s = %s\n", inst.Result_addr, inst.Var_um)
				} else {
					fmt.Printf("│       %s = %s %s %s\n", inst.Result_addr, inst.Var_um, inst.Operador, inst.Var_dois)
				}
			}
			// --- INTEGRAÇÃO DO CÓDIGO DESTINO (PYTHON) ---
			fmt.Println("│\n│   [Python] Código Destino gerado:")
			codigoPython := compiler.GerarPython(gerador.Instruction_list)

			// Quebra o texto gerado em linhas para alinhar com o design do seu terminal
			linhasPython := strings.Split(strings.TrimSpace(codigoPython), "\n")
			for _, linha := range linhasPython {
				fmt.Printf("│       %s\n", linha)
			}
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
