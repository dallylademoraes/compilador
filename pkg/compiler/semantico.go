package compiler

import "fmt"

// AnalisadorSemantico percorre a AST e valida o uso de identificadores.
type AnalisadorSemantico struct {
	tabela *TabelaSimbolos
	erros  []string
}

// NovoAnalisadorSemantico cria uma instância pronta para validar uma AST.
func NovoAnalisadorSemantico() *AnalisadorSemantico {
	return &AnalisadorSemantico{
		tabela: NovaTabelaSimbolos(),
	}
}

// AnalisarSemantica executa a Fase 4 sobre a AST informada.
func AnalisarSemantica(raiz *Program) (*TabelaSimbolos, []string) {
	analisador := NovoAnalisadorSemantico()
	analisador.Check(raiz)
	return analisador.Tabela(), analisador.Erros()
}

// Check executa a validação semântica completa.
func (a *AnalisadorSemantico) Check(raiz *Program) bool {
	if raiz == nil {
		a.adicionarErro("[Erro semântico] AST vazia")
		return false
	}

	raiz.Accept(a)

	return len(a.erros) == 0
}

// Tabela devolve a tabela de símbolos construída.
func (a *AnalisadorSemantico) Tabela() *TabelaSimbolos {
	return a.tabela
}

// Erros devolve uma cópia dos erros acumulados.
func (a *AnalisadorSemantico) Erros() []string {
	erros := make([]string, len(a.erros))
	copy(erros, a.erros)
	return erros
}

// VisitProgram percorre a raiz da AST.
func (a *AnalisadorSemantico) VisitProgram(no *Program) {
	for _, stmt := range no.Statements {
		stmt.Accept(a)
	}
}

// VisitDecl registra uma declaração e valida sua expressão inicial, se existir.
func (a *AnalisadorSemantico) VisitDecl(no *Decl) {
	if no.Value != nil {
		no.Value.Accept(a)
	}

	valor := ""
	if no.Value != nil {
		valor = formatarExpr(no.Value)
	}

	if err := a.tabela.InserirComValor(no.Name, "int", valor); err != nil {
		a.adicionarErro(err.Error())
	}
}

// VisitAssign valida a expressão e atualiza um símbolo já declarado.
func (a *AnalisadorSemantico) VisitAssign(no *Assign) {
	no.Value.Accept(a)

	if err := a.tabela.Atribuir(no.Name, formatarExpr(no.Value)); err != nil {
		a.adicionarErro(err.Error())
	}
}

// VisitBinOp percorre recursivamente os operandos de uma expressão binária.
func (a *AnalisadorSemantico) VisitBinOp(no *BinOp) {
	no.Left.Accept(a)
	no.Right.Accept(a)
}

// VisitIdentifier valida o uso de um identificador em expressão.
func (a *AnalisadorSemantico) VisitIdentifier(no *Identifier) {
	if _, ok := a.tabela.Buscar(no.Name); !ok {
		a.adicionarErro(fmt.Sprintf("[Erro semântico] Variável '%s' usada antes de ser declarada", no.Name))
	}
}

// VisitNum não exige validações semânticas adicionais.
func (a *AnalisadorSemantico) VisitNum(*Num) {}

func (a *AnalisadorSemantico) adicionarErro(erro string) {
	a.erros = append(a.erros, erro)
}

func formatarExpr(expr Expr) string {
	switch no := expr.(type) {
	case *Identifier:
		return no.Name
	case *Num:
		return no.Value
	case *BinOp:
		return fmt.Sprintf("(%s %s %s)", formatarExpr(no.Left), no.Op, formatarExpr(no.Right))
	default:
		return ""
	}
}
