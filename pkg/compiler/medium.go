package compiler

import "fmt"


type Instruction struct{
	Var_um string
	Var_dois string
	Operador string
	Result_addr string
}

type Gerador3AC struct{
	Instruction_list []Instruction
	Counter int
}

func (g *Gerador3AC) newTemp() string{
	temp := fmt.Sprintf("t%d", g.Counter)
	g.Counter++
	return temp
}

func (g *Gerador3AC)VisitNode(no Node) string{
	switch n := no.(type){
		case *Num:
			return string(n.Value)
		case *Identifier:
			return string(n.Name)
		case *BinOp:
			return g.binOpFunc(n)
		case *Assign:
			return g.assignFunc(n)
		case *Decl:
			return g.declFunc(n)
		default:
			panic("NODE NOT RECOGNIZABLE")
	}
}




func (g *Gerador3AC) declFunc(no *Decl) string{
	right_value := g.VisitNode(no.Value)
	if no.Value != nil{
		new_instruction := Instruction{
			Var_um: right_value,
			Var_dois: "",
			Operador: "=",
			Result_addr: no.Name,
		}

		g.Instruction_list = append(g.Instruction_list, new_instruction)
	}
	return no.Name
}

func (g *Gerador3AC) assignFunc(no *Assign) string{
	right_value := g.VisitNode(no.Value)

	new_instruction := Instruction{
		Var_um: right_value,
		Var_dois: "",
		Operador: "=",
		Result_addr: no.Name,
	}
	g.Instruction_list = append(g.Instruction_list, new_instruction)
	return no.Name
}

func(g *Gerador3AC) binOpFunc(no *BinOp) string{
	left_node := g.VisitNode(no.Left)
	right_node := g.VisitNode(no.Right)
	result_addr := g.newTemp()

	new_instruction := Instruction{
		Var_um: left_node,
		Var_dois: right_node,
		Operador: no.Op,
		Result_addr: result_addr,
	}
	g.Instruction_list = append(g.Instruction_list, new_instruction)
	return result_addr
}
