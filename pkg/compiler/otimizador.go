package compiler

import (
	"fmt"
	"strconv"
)

func Otimizar(instrs []Instruction) []Instruction {
	instrs = constantFolding(instrs)
	instrs = constantPropagation(instrs)
	instrs = deadCodeElimination(instrs)
	return instrs
}

func constantFolding(instrs []Instruction) []Instruction {
	resultado := make([]Instruction, 0, len(instrs))

	for _, inst := range instrs {
		if inst.Operador != "=" && inst.Var_um != "" && inst.Var_dois != "" {
			esq, errEsq := strconv.Atoi(inst.Var_um)
			dir, errDir := strconv.Atoi(inst.Var_dois)

			if errEsq == nil && errDir == nil {
				valor, ok := calcular(esq, dir, inst.Operador)
				if ok {
					inst = Instruction{
						Result_addr: inst.Result_addr,
						Operador:    "=",
						Var_um:      fmt.Sprintf("%d", valor),
						Var_dois:    "",
					}
				}
			}
		}
		resultado = append(resultado, inst)
	}

	return resultado
}

func calcular(esq, dir int, op string) (int, bool) {
	switch op {
	case "+":
		return esq + dir, true
	case "-":
		return esq - dir, true
	case "*":
		return esq * dir, true
	case "/":
		if dir == 0 {
			return 0, false
		}
		return esq / dir, true
	}
	return 0, false
}

func constantPropagation(instrs []Instruction) []Instruction {
	constantes := make(map[string]string)

	resultado := make([]Instruction, 0, len(instrs))

	for _, inst := range instrs {
		if val, ok := constantes[inst.Var_um]; ok {
			inst.Var_um = val
		}
		if val, ok := constantes[inst.Var_dois]; ok {
			inst.Var_dois = val
		}

		if inst.Operador != "=" && inst.Var_um != "" && inst.Var_dois != "" {
			esq, errEsq := strconv.Atoi(inst.Var_um)
			dir, errDir := strconv.Atoi(inst.Var_dois)
			if errEsq == nil && errDir == nil {
				if valor, ok := calcular(esq, dir, inst.Operador); ok {
					inst = Instruction{
						Result_addr: inst.Result_addr,
						Operador:    "=",
						Var_um:      fmt.Sprintf("%d", valor),
						Var_dois:    "",
					}
				}
			}
		}

		if inst.Operador == "=" && inst.Var_dois == "" {
			if _, err := strconv.Atoi(inst.Var_um); err == nil {
				constantes[inst.Result_addr] = inst.Var_um
			}
		}

		resultado = append(resultado, inst)
	}

	return resultado
}

func deadCodeElimination(instrs []Instruction) []Instruction {
	usos := make(map[string]int)
	for _, inst := range instrs {
		if inst.Var_um != "" {
			usos[inst.Var_um]++
		}
		if inst.Var_dois != "" {
			usos[inst.Var_dois]++
		}
	}

	resultado := make([]Instruction, 0, len(instrs))
	for _, inst := range instrs {
		if eTemporaria(inst.Result_addr) && usos[inst.Result_addr] == 0 {
			continue
		}
		resultado = append(resultado, inst)
	}

	return resultado
}

func eTemporaria(nome string) bool {
	if len(nome) < 2 || nome[0] != 't' {
		return false
	}
	_, err := strconv.Atoi(nome[1:])
	return err == nil
}