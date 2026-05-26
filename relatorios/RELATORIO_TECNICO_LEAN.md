# 📝 Relatório Técnico: MiniCompilador (Foco: Lean — 3AC & Python)

Este relatório consolida o estado atual do projeto **MiniCompilador** para subsidiar a implementação das fases de **Código Intermediário (3AC)** e **Geração de Código Python**.

---

## 1. Visão Geral do Sistema
O compilador está estruturado em Go, utilizando `goyacc` para o parser LALR(1). O fluxo atual cobre as fases de **Análise Léxica**, **Análise Sintática** e **Construção da AST**.

### Arquitetura do Pipeline
1.  **Lexer (`lexer.go`):** Tokenização da entrada.
2.  **Parser (`parser.y`):** Validação gramatical e construção da AST.
3.  **AST (`ast.go`):** Representação em memória e exportação visual (DOT/PNG).
4.  **Tabela de Símbolos (`semantico.go`):** *Em desenvolvimento (Ícaro).*
5.  **Código Intermediário (`intermediario.go`):** **SUA TAREFA (Lean).**
6.  **Geração de Código Python (`gerador.go`):** **SUA TAREFA (Lean).**

---

## 2. O que já está PRONTO (Fase João/Frontend)

### A Gramática (BNF)
O parser já reconhece e valida:
- Declarações (`int x;`, `int x = 10;`)
- Atribuições (`x = a + b;`)
- Operações binárias (`+`, `-`, `*`, `/`)
- Precedência de operadores (parênteses e `*`/`/` sobre `+`/`-`)

### A Estrutura da AST (`pkg/compiler/ast.go`)
Você consumirá estes nós principais:
- `*Program`: Contém `Statements []Stmt`.
- `*Decl`: Possui `Name string` e `Value Expr` (opcional).
- `*Assign`: Possui `Name string` e `Value Expr`.
- `*BinOp`: Possui `Op string`, `Left Expr` e `Right Expr`.
- `*Identifier` / `*Num`: Folhas da árvore.

### Validação
- O arquivo `cmd/minicompiler/main_test.go` contém 9 casos de teste que já passam, cobrindo expressões complexas e divisões.
- O executável gera visualizações automáticas em `ast/ast.png`.

---

## 3. O que VOCÊ (Lean) precisa implementar

### Parte A: Código Intermediário (3AC) — `intermediario.go`
Você deve transformar a AST em instruções de três endereços. 
**Exemplo de entrada:** `int custo = tempo * usuarios;`
**Exemplo de 3AC esperado:**
```text
t0 = tempo * usuarios
custo = t0
```

**Requisitos técnicos:**
1.  **Gerador de Temporárias:** Função que retorna `t0, t1, t2...` a cada chamada.
2.  **Estrutura de Instrução:**
    ```go
    type Instrucao struct {
        Op, Arg1, Arg2, Result string
    }
    ```
3.  **Percorrimento (Visitor):** Uma função recursiva que desce na AST. Ao encontrar um `BinOp`, ela gera as instruções para os filhos, cria um novo temporário para o resultado e retorna esse temporário para o nó pai.

### Parte B: Geração de Código Python — `gerador.go`
Traduzir as instruções 3AC (ou a AST diretamente) para Python.
- `int x = 10;` ➔ `x = 10`
- `int tempo = p + r;` ➔ `tempo = p + r`
- Python não exige declaração de tipos, o que simplifica sua vida.

---

## 4. Arquivos-Chave para Estudo Imediato
1.  `pkg/compiler/ast.go`: Entenda como os nós estão estruturados.
2.  `pkg/compiler/parser.y`: Veja como a árvore é montada.
3.  `script_teste_prof.txt`: O código de teste oficial que seu gerador deve processar.

---

## 5. Roadmap Sugerido
1.  Criar `pkg/compiler/ircode.go` para definir a estrutura `Instruction` e o gerador de temporários.
2.  Criar a lógica de percorrer a lista de `Statements` da AST no `main.go`.
3.  Integrar a saída do seu código como a "Fase 5" e "Fase 6" no console.

---
*Relatório gerado em 25/05/2026 para a equipe do MiniCompilador.*
