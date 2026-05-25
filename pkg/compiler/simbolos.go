package compiler

import (
	"fmt"
	"sort"
	"strings"
)

// Simbolo representa uma entrada da tabela de símbolos.
type Simbolo struct {
	Nome  string
	Tipo  string
	Valor string
}

// TabelaSimbolos armazena os identificadores declarados no programa.
type TabelaSimbolos struct {
	tabela map[string]Simbolo
}

// NovaTabelaSimbolos cria uma tabela de símbolos vazia.
func NovaTabelaSimbolos() *TabelaSimbolos {
	return &TabelaSimbolos{
		tabela: make(map[string]Simbolo),
	}
}

// Inserir registra um novo símbolo na tabela.
func (t *TabelaSimbolos) Inserir(nome, tipo string) error {
	return t.InserirComValor(nome, tipo, "")
}

// InserirComValor registra um novo símbolo com valor inicial.
func (t *TabelaSimbolos) InserirComValor(nome, tipo, valor string) error {
	if _, existe := t.tabela[nome]; existe {
		return fmt.Errorf("[Erro semântico] Variável '%s' já foi declarada", nome)
	}

	t.tabela[nome] = Simbolo{
		Nome:  nome,
		Tipo:  tipo,
		Valor: valor,
	}

	return nil
}

// Atribuir atualiza o valor de um símbolo existente.
func (t *TabelaSimbolos) Atribuir(nome, valor string) error {
	simbolo, existe := t.tabela[nome]
	if !existe {
		return fmt.Errorf("[Erro semântico] Variável '%s' usada antes de ser declarada", nome)
	}

	simbolo.Valor = valor
	t.tabela[nome] = simbolo
	return nil
}

// Buscar procura um símbolo pelo nome.
func (t *TabelaSimbolos) Buscar(nome string) (Simbolo, bool) {
	simbolo, ok := t.tabela[nome]
	return simbolo, ok
}

// Imprimir exibe a tabela formatada no terminal.
func (t *TabelaSimbolos) Imprimir() {
	fmt.Print(t.String())
}

// String retorna a tabela formatada.
func (t *TabelaSimbolos) String() string {
	var sb strings.Builder

	sb.WriteString("┌─────────────────────┬─────────────────────┬─────────────────────┐\n")
	sb.WriteString(fmt.Sprintf("│ %-19s │ %-19s │ %-19s │\n", "Nome", "Tipo", "Valor"))
	sb.WriteString("├─────────────────────┼─────────────────────┼─────────────────────┤\n")

	for _, simbolo := range t.ordenados() {
		valor := simbolo.Valor
		if valor == "" {
			valor = "<não inicializado>"
		}
		sb.WriteString(fmt.Sprintf("│ %-19s │ %-19s │ %-19s │\n", simbolo.Nome, simbolo.Tipo, valor))
	}

	sb.WriteString("└─────────────────────┴─────────────────────┴─────────────────────┘\n")
	return sb.String()
}

func (t *TabelaSimbolos) ordenados() []Simbolo {
	nomes := make([]string, 0, len(t.tabela))
	for nome := range t.tabela {
		nomes = append(nomes, nome)
	}
	sort.Strings(nomes)

	simbolos := make([]Simbolo, 0, len(nomes))
	for _, nome := range nomes {
		simbolos = append(simbolos, t.tabela[nome])
	}

	return simbolos
}
