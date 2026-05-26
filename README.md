# 🛠️ MiniCompilador — Disciplina de Compiladores

> Linguagem: **Go** · Parser: **goyacc (LALR(1))**  
> Repositório: `compilador/pkg/compiler`  
> Grupo: **João Walcacer · Ícaro · Lean · Dallyla**

---

## 📂 Estrutura do Projeto

| Diretório / Arquivo | Função |
|---|---|
| `cmd/minicompiler/main.go` | Ponto de entrada da aplicação (CLI) |
| `cmd/minicompiler/main_test.go` | Conjunto de testes e exemplos de validação |
| `pkg/compiler/ast.go` | Estruturas da AST + impressão + exportação DOT/PNG |
| `pkg/compiler/lexer.go` | Analisador léxico — identifica tokens |
| `pkg/compiler/parser.y` | Definição da gramática (Yacc) |
| `pkg/compiler/parser.go` | Código gerado pelo goyacc (não editar diretamente) |
| `ast/` | Saída dos artefatos da árvore (`ast.dot` e `ast.png`) |
| `Makefile` | Atalhos para build, testes e geração de código |

---

## 📌 Estado atual do projeto

| Arquivo | O que faz | Status |
|---|---|---|
| `lexer.go` | Análise léxica — tokeniza `int`, IDs, NUMs, operadores (`+`, `-`, `*`, `/`) e delimitadores | ✅ Pronto |
| `parser.y` | Gramática LALR(1) via goyacc — valida declarações e expressões com `+`, `-`, `*`, `/` | ✅ Pronto |
| `parser.go` | Código gerado pelo goyacc | ✅ Gerado |
| `main.go` | Entrada do programa — lê do terminal (ou arquivo), fases léxica/sintática/AST e exportação em `ast/` | ✅ Pronto |
| `main_test.go` | 9 casos de teste (aceitos e rejeitados), incluindo divisão | ✅ Pronto |
| `ast.go` | Árvore sintática abstrata + impressão hierárquica + exportação DOT/PNG | ✅ Pronto |
| `simbolos.go` + `semantico.go` | Tabela de símbolos + análise semântica | ⏳ Ícaro |
| `intermediario.go` + `gerador.go` | Código intermediário + geração de código final | ⏳ Lean |
| `otimizador.go` | Otimizações + relatório + apresentação | ⏳ Dallyla |

---

## 🚀 Como Executar

### Usando o Makefile (recomendado)

```bash
make test      # Rodar os testes
make run       # Modo interativo
make build     # Gerar o executável
```

### Usando Go CLI

```bash
go test -v ./...              # Testes
go run ./cmd/minicompiler      # Modo interativo
```

### Alterando a gramática

Se modificar `parser.y`, regere o parser:

```bash
go install golang.org/x/tools/cmd/goyacc@latest
make generate
```

---

## 📐 Gramática da MiniLang (BNF)

```
programa   → lista_stmt
lista_stmt → lista_stmt stmt | stmt
stmt       → int id ;
           | int id = expr ;
           | id = expr ;
expr       → expr + termo | expr - termo | termo
termo      → termo * fator | termo / fator | fator
fator      → id | num | ( expr )
```

**Tokens reconhecidos:** `int` · identificadores · números inteiros · `=` · `+` · `-` · `*` · `/` · `;` · `(` · `)`

### Exemplos de código aceitos

```text
int x = 10;
total = (a + b) * 5;
int resultado = x + y + z;
```

---

## 📝 Pipeline de Execução

```
Código Fonte
    ↓
Lexer (lexer.go)          → tokens
    ↓
Parser (parser.go)        → validação sintática
    ↓
AST (ast.go)              → árvore de nós em memória        ← João
    ↓
Tabela de Símbolos        → validação semântica             ← Ícaro
(simbolos.go / semantico.go)
    ↓
Código Intermediário      → instruções de três endereços    ← Lean
(intermediario.go)
    ↓
Geração de Código         → código Python válido            ← Lean
(gerador.go)
    ↓
Otimizações               → propagação de constantes,       ← Dallyla
(otimizador.go)             eliminação de código morto
```

---
