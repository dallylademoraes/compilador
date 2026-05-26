package compiler

import (
	"fmt"
	"strings"
)

func GerarPython(instructions []Instruction) string{

	var new_string strings.Builder
	new_string.WriteString("# Script gerado automaticamente pelo compilador\n\n")
	for _, instr := range instructions{
		if instr.Operador == "="{
			line := fmt.Sprintf("%s = %s\n", instr.Result_addr, instr.Var_um)
			new_string.WriteString(line)
		}else{
			line := fmt.Sprintf("%s = %s %s %s\n", instr.Result_addr, instr.Var_um, instr.Operador, instr.Var_dois)
			new_string.WriteString(line)
		}
	}
	return new_string.String()
}
