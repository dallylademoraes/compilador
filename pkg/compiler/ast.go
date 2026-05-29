package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Node representa qualquer nó da AST.
type Node interface {
	node()
	Accept(Visitor)
}

// Visitor define operações sobre os nós da AST.
type Visitor interface {
	VisitProgram(*Program)
	VisitDecl(*Decl)
	VisitAssign(*Assign)
	VisitBinOp(*BinOp)
	VisitIdentifier(*Identifier)
	VisitNum(*Num)
}

// Stmt representa um statement da linguagem.
type Stmt interface {
	Node
	stmt()
}

// Expr representa uma expressão da linguagem.
type Expr interface {
	Node
	expr()
}

// Program é a raiz da AST.
type Program struct {
	Statements []Stmt
}

// Decl representa declaração de variável: int x; ou int x = expr;
type Decl struct {
	Name  string
	Value Expr
}

// Assign representa atribuição: x = expr;
type Assign struct {
	Name  string
	Value Expr
}

// BinOp representa operação binária.
type BinOp struct {
	Op    string
	Left  Expr
	Right Expr
}

// Identifier representa referência a identificador.
type Identifier struct {
	Name string
}

// Num representa literal numérico.
type Num struct {
	Value string
}

func (*Program) node()    {}
func (*Decl) node()       {}
func (*Decl) stmt()       {}
func (*Assign) node()     {}
func (*Assign) stmt()     {}
func (*BinOp) node()      {}
func (*BinOp) expr()      {}
func (*Identifier) node() {}
func (*Identifier) expr() {}
func (*Num) node()        {}
func (*Num) expr()        {}

func (n *Program) Accept(v Visitor)    { v.VisitProgram(n) }
func (n *Decl) Accept(v Visitor)       { v.VisitDecl(n) }
func (n *Assign) Accept(v Visitor)     { v.VisitAssign(n) }
func (n *BinOp) Accept(v Visitor)      { v.VisitBinOp(n) }
func (n *Identifier) Accept(v Visitor) { v.VisitIdentifier(n) }
func (n *Num) Accept(v Visitor)        { v.VisitNum(n) }

// ImprimirAST exibe a árvore sintática em formato hierárquico.
func ImprimirAST(no Node, indent int) {
	if no == nil {
		fmt.Println("<AST vazia>")
		return
	}

	prefix := strings.Repeat(" ", indent)
	printRoot(no, prefix)
}

func printRoot(no Node, prefix string) {
	switch n := no.(type) {
	case *Program:
		fmt.Printf("%sPrograma\n", prefix)
		for i, s := range n.Statements {
			printNode(s, prefix, i == len(n.Statements)-1)
		}
	default:
		printNode(n, prefix, true)
	}
}

func printNode(no Node, prefix string, last bool) {
	branch := "├── "
	nextPrefix := prefix + "│   "
	if last {
		branch = "└── "
		nextPrefix = prefix + "    "
	}

	switch n := no.(type) {
	case *Decl:
		fmt.Printf("%s%sDecl: int %s\n", prefix, branch, n.Name)
		if n.Value != nil {
			printNode(n.Value, nextPrefix, true)
		}
	case *Assign:
		fmt.Printf("%s%sAtrib: %s\n", prefix, branch, n.Name)
		printNode(n.Value, nextPrefix, true)
	case *BinOp:
		fmt.Printf("%s%sBinOp: %s\n", prefix, branch, n.Op)
		printNode(n.Left, nextPrefix, false)
		printNode(n.Right, nextPrefix, true)
	case *Identifier:
		fmt.Printf("%s%sID: %s\n", prefix, branch, n.Name)
	case *Num:
		fmt.Printf("%s%sNum: %s\n", prefix, branch, n.Value)
	case *Program:
		fmt.Printf("%s%sPrograma\n", prefix, branch)
		for i, s := range n.Statements {
			printNode(s, nextPrefix, i == len(n.Statements)-1)
		}
	default:
		fmt.Printf("%s%s<no desconhecido>\n", prefix, branch)
	}
}

// ExportarASTDOT salva a AST no formato Graphviz DOT.
func ExportarASTDOT(no Node, caminho string) error {
	if no == nil {
		return fmt.Errorf("AST vazia")
	}

	var sb strings.Builder
	sb.WriteString("digraph AST {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=box, style=rounded, fontname=\"Arial\"];\n")

	nextID := 0
	var walk func(Node) int
	walk = func(n Node) int {
		id := nextID
		nextID++

		label := rotuloNo(n)
		sb.WriteString(fmt.Sprintf("  n%d [label=%q];\n", id, label))

		for _, filho := range filhosNo(n) {
			childID := walk(filho)
			sb.WriteString(fmt.Sprintf("  n%d -> n%d;\n", id, childID))
		}

		return id
	}

	walk(no)
	sb.WriteString("}\n")

	return os.WriteFile(caminho, []byte(sb.String()), 0o644)
}

// ExportarASTPNG salva a AST em PNG usando o comando `dot` do Graphviz.
// Também salva o arquivo .dot no mesmo diretório para inspeção.
func ExportarASTPNG(no Node, caminhoPNG string) error {
	if no == nil {
		return fmt.Errorf("AST vazia")
	}

	dotPath := strings.TrimSuffix(caminhoPNG, filepath.Ext(caminhoPNG)) + ".dot"
	if err := ExportarASTDOT(no, dotPath); err != nil {
		return err
	}

	dotCmd, err := resolverComandoDot()
	if err != nil {
		return err
	}

	cmd := exec.Command(dotCmd, "-Tpng", dotPath, "-o", caminhoPNG)
	if out, err := cmd.CombinedOutput(); err != nil {
		if len(out) > 0 {
			return fmt.Errorf("falha ao executar dot: %v (%s)", err, strings.TrimSpace(string(out)))
		}
		return fmt.Errorf("falha ao executar dot: %v", err)
	}

	return nil
}

func resolverComandoDot() (string, error) {
	// Permite sobrescrever o caminho do Graphviz explicitamente.
	if dotEnv := strings.TrimSpace(os.Getenv("GRAPHVIZ_DOT")); dotEnv != "" {
		if _, err := os.Stat(dotEnv); err == nil {
			return dotEnv, nil
		}
	}

	if caminho, err := exec.LookPath("dot"); err == nil {
		return caminho, nil
	}

	if runtime.GOOS == "windows" {
		fallbacks := []string{
			`C:\Program Files\Graphviz\bin\dot.exe`,
			`C:\Program Files (x86)\Graphviz\bin\dot.exe`,
		}
		for _, caminho := range fallbacks {
			if _, err := os.Stat(caminho); err == nil {
				return caminho, nil
			}
		}
	}

	return "", fmt.Errorf("falha ao executar dot: executable file not found in %%PATH%%")
}

func rotuloNo(no Node) string {
	switch n := no.(type) {
	case *Program:
		return "Programa"
	case *Decl:
		return fmt.Sprintf("Decl: int %s", n.Name)
	case *Assign:
		return fmt.Sprintf("Atrib: %s", n.Name)
	case *BinOp:
		return fmt.Sprintf("BinOp: %s", n.Op)
	case *Identifier:
		return fmt.Sprintf("ID: %s", n.Name)
	case *Num:
		return fmt.Sprintf("Num: %s", n.Value)
	default:
		return "<no desconhecido>"
	}
}

func filhosNo(no Node) []Node {
	switch n := no.(type) {
	case *Program:
		filhos := make([]Node, 0, len(n.Statements))
		for _, s := range n.Statements {
			filhos = append(filhos, s)
		}
		return filhos
	case *Decl:
		if n.Value == nil {
			return nil
		}
		return []Node{n.Value}
	case *Assign:
		return []Node{n.Value}
	case *BinOp:
		return []Node{n.Left, n.Right}
	default:
		return nil
	}
}
