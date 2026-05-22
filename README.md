# 🛠️ MiniCompilador — Divisão de Tarefas

> Disciplina de Compiladores · Grupo de 4 pessoas  
> Linguagem: **Go** · Parser: **goyacc (LALR(1))**  
> Repositório: `compilador/pkg/compiler`

---

## 📌 Estado atual do projeto

| Arquivo | O que faz | Status |
|---|---|---|
| `lexer.go` | Análise léxica — tokeniza `int`, IDs, NUMs, operadores, delimitadores | ✅ Pronto |
| `parser.y` | Gramática LALR(1) via goyacc — valida declarações e expressões | ✅ Pronto |
| `parser.go` | Código gerado pelo goyacc (não editar diretamente) | ✅ Gerado |
| `main.go` | Entrada do programa — modo interativo `-i` + saída formatada | ✅ Pronto |
| `main_test.go` | 8 casos de teste (aceitos e rejeitados) | ✅ Pronto |

**O que ainda falta:** AST estruturada, tabela de símbolos, código intermediário, geração de código final, relatório e apresentação.

---

## 👥 Divisão por Pessoa

---

### 👤 Pessoa A — AST (Árvore Sintática Abstrata)
> Corresponde às **Seções 4 e 5** do trabalho (Análise Sintática + Tradução Dirigida por Sintaxe)

**Contexto:** Hoje o parser só imprime mensagens no console nas ações semânticas (`fmt.Printf`). Isso precisa evoluir para construir uma árvore de nós em memória que represente o programa.

**Tarefas:**

1. Criar o arquivo `ast.go` com as structs de nó:
   - `NóPrograma` — raiz da árvore, contém lista de statements
   - `NóDecl` — declaração `int x` ou `int x = expr`
   - `NóAtrib` — atribuição `x = expr`
   - `NóBinOp` — operação binária com campos `Op`, `Esquerda`, `Direita`
   - `NóID` — referência a variável (folha da árvore)
   - `NóNum` — número literal (folha da árvore)

2. Modificar as **ações semânticas do `parser.y`** para, em vez de só imprimir, popular a AST construindo nós e conectando-os.

3. Implementar uma função `ImprimirAST(nó Nó, indent int)` que exibe a árvore de forma visual no terminal, como:
   ```
   Programa
   └── Decl: int soma
       └── BinOp: +
           ├── ID: a
           └── ID: b
   ```

4. Adicionar a exibição da AST como **Fase 3** no `main.go` e nos testes.

**Entrega para o grupo:** a AST populada deve ser retornada pela função `ExecutarParser` para que a Pessoa C possa percorrê-la.

---

### 👤 Pessoa B — Tabela de Símbolos + Análise Semântica
> Corresponde às **Seções 5 e 7** do trabalho (Tradução Dirigida por Sintaxe + Ambientes de Execução)

**Contexto:** A linguagem tem variáveis tipadas (`int`). É necessário rastrear quais variáveis foram declaradas, seus tipos e valores atuais para detectar erros semânticos.

**Tarefas:**

1. Criar o arquivo `simbolos.go` com a struct `TabelaSimbolos`:
   ```go
   type Simbolo struct {
       Nome  string
       Tipo  string // "int"
       Valor string // pode ser "" se ainda não inicializado
   }
   type TabelaSimbolos struct {
       tabela map[string]Simbolo
   }
   ```

2. Implementar os métodos:
   - `Inserir(nome, tipo string)` — registra nova variável; erro se já declarada
   - `Atribuir(nome, valor string)` — atualiza valor; erro se não declarada
   - `Buscar(nome string) (Simbolo, bool)` — consulta variável
   - `Imprimir()` — exibe a tabela formatada no terminal

3. Criar um **visitor semântico** (`semantico.go`) que percorre a AST (produzida pela Pessoa A) e alimenta a tabela de símbolos:
   - Ao visitar `NóDecl`: chama `Inserir`
   - Ao visitar `NóAtrib`: verifica se a variável existe antes de aceitar
   - Ao visitar `NóID` dentro de uma expressão: verifica se foi declarado

4. Reportar erros semânticos com mensagens claras:
   ```
   [Erro semântico] Variável 'x' usada antes de ser declarada
   [Erro semântico] Variável 'y' já foi declarada
   ```

5. Adicionar a exibição da tabela de símbolos como **Fase 4** no `main.go`.

**Entrega para o grupo:** a tabela de símbolos preenchida deve estar disponível para a Pessoa C usar na geração de código.

---

### 👤 Pessoa C — Código Intermediário + Geração de Código Final
> Corresponde às **Seções 6 e 8** do trabalho (Geração de Código Intermediário + Geração de Código)

**Contexto:** Com a AST pronta (Pessoa A) e a tabela de símbolos validada (Pessoa B), esta fase transforma o programa em código executável em dois passos.

**Tarefas — Parte 1: Código Intermediário (`intermediario.go`)**

1. Implementar um gerador de **código de três endereços** que percorre a AST:
   - Cada expressão binária gera um temporário: `t1 = a + b`
   - Declarações viram: `int soma` seguido de `soma = t1`
   - Atribuições simples: `x = 42`

2. Exemplo esperado para `int soma = a + b * c;`:
   ```
   t1 = b * c
   t2 = a + t1
   soma = t2
   ```

3. Armazenar as instruções como lista de structs:
   ```go
   type Instrucao struct {
       Resultado string
       Op        string
       Arg1      string
       Arg2      string
   }
   ```

4. Exibir o código intermediário como **Fase 5** no `main.go`.

**Tarefas — Parte 2: Geração de Código Final (`gerador.go`)**

1. Traduzir a lista de instruções de três endereços para **código Python válido**:
   - `int x;` → `x = None`
   - `int soma = a + b;` → `soma = a + b`
   - `x = 42;` → `x = 42`

2. Salvar o arquivo `.py` gerado quando rodado com flag `-o`:
   ```
   go run . -i -o saida.py
   ```

3. Exibir o código gerado como **Fase 6** no `main.go`.

**Entrega para o grupo:** o código intermediário e final devem funcionar com qualquer entrada que passe pela validação semântica da Pessoa B.

---

### 👤 Pessoa D — Otimizações + Relatório + Apresentação
> Corresponde às **Seções 9, 10, 11 e 12** do trabalho

**Contexto:** As otimizações podem ser aplicadas sobre a lista de instruções de três endereços gerada pela Pessoa C. O relatório deve cobrir todas as seções do trabalho.

**Tarefas — Implementação (`otimizador.go`)**

1. Implementar **propagação de constantes**: se `t1 = 5 + 3`, substituir por `t1 = 8` diretamente, sem emitir a instrução de soma.

2. Implementar **eliminação de código morto**: se um temporário `t1` é calculado mas nunca usado em nenhuma outra instrução, removê-lo da lista.

3. Exibir o código intermediário **antes e depois** da otimização para evidenciar a diferença.

**Tarefas — Relatório**

Escrever o relatório final cobrindo todas as 12 seções do trabalho. Para cada seção:

| Seção | O que escrever |
|---|---|
| 1 – Introdução | Descrição da MiniLang, exemplos de programas válidos |
| 2 – Gramática e Autômatos | BNF do `parser.y` + diagrama do autômato léxico |
| 3 – Análise Léxica | Explicar o `lexer.go`, mostrar tabela de tokens de exemplo |
| 4 – Análise Sintática | Explicar o parser LALR(1), mostrar exemplo de AST gerada |
| 5 – Tradução Dirigida | Explicar as ações semânticas e a tabela de símbolos |
| 6 – Código Intermediário | Mostrar exemplos de três endereços gerados |
| 7 – Ambiente de Execução | Descrever como a tabela de símbolos gerencia memória |
| 8 – Geração de Código | Mostrar exemplo de código Python gerado |
| 9 – Otimizações | Mostrar antes/depois da propagação de constantes e eliminação de mortos |
| 10 – Paralelismo | Analisar teoricamente quais instruções do código intermediário são independentes |
| 11 – Localidade | Discutir como reordenar instruções poderia melhorar cache |
| 12 – Aplicações | Relacionar com JIT (Java/.NET), compiladores para GPU, sistemas embarcados |

**Tarefas — Apresentação**

Montar slides com no mínimo:
- Slide 1: Título e integrantes
- Slide 2: Visão geral da arquitetura do compilador
- Slides 3–8: Uma fase por slide (léxico, sintático, AST, tabela, intermediário, código final)
- Slide 9: Otimizações com exemplo
- Slide 10: Demo ao vivo + conclusão

---

## 🔗 Pontos de integração entre as pessoas

```
Pessoa A (AST) ──────────────────────────────────────────→ Pessoa C (usa a AST pra gerar código)
                                                         ↗
Pessoa B (Tabela de Símbolos) ───────────────────────────

Pessoa C (Código Intermediário) ─────────────────────────→ Pessoa D (aplica otimizações)
```

| Handoff | Quem entrega | Quem recebe | O que precisa estar pronto |
|---|---|---|---|
| `NóPrograma` populado + `ImprimirAST` | Pessoa A | Pessoa C | Structs de nó definidas, `parser.y` gerando a árvore |
| `TabelaSimbolos` + `Visitor` semântico | Pessoa B | Pessoa C | Tabela funcionando, erros semânticos reportados |
| Lista de `Instrucao` (três endereços) | Pessoa C | Pessoa D | Gerador percorrendo a AST corretamente |

---

## ▶️ Como rodar o projeto

```bash
# Gerar o parser (necessário após editar parser.y)
goyacc -o pkg/compiler/parser.go -p yy pkg/compiler/parser.y

# Rodar em modo interativo
go run . -i

# Rodar os testes
go test ./...
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
termo      → termo * fator | fator
fator      → id | num | ( expr )
```

**Tokens reconhecidos:** `int` · identificadores · números inteiros · `=` · `+` · `-` · `*` · `;` · `(` · `)`

# Minicompiler – Disciplina de Compiladores

Este projeto é um compilador didático desenvolvido em Go, focado nas fases de **Análise Léxica** e **Análise Sintática**. Ele utiliza `goyacc` para o parsing LALR(1).

---

## 📂 Estrutura do Projeto

A estrutura segue as boas práticas de organização do ecossistema Go:

| Diretório / Arquivo | Função |
|--- |--- |
| `cmd/minicompiler/main.go` | Ponto de entrada da aplicação (CLI). |
| `cmd/minicompiler/main_test.go` | Conjunto de testes e exemplos de validação. |
| `pkg/compiler/` | O motor do compilador (Lexer e Parser). |
| `pkg/compiler/lexer.go` | Analisador léxico que identifica tokens. |
| `pkg/compiler/parser.y` | Definição da gramática (Yacc). |
| `Makefile` | Atalhos para automação de build, testes e geração de código. |

---

## 🚀 Como Executar

### 1. Pré-requisitos
Certifique-se de ter o **Go** instalado (versão 1.22 ou superior recomendada).

### 2. Usando o Makefile (Recomendado)
Se você estiver no Linux ou Mac, use os comandos simplificados:

*   **Validar exemplos (Testes):**
    ```bash
    make test
    ```
*   **Rodar modo interativo (Testar seu código):**
    ```bash
    make run
    ```
*   **Gerar o executável:**
    ```bash
    make build
    ```

### 3. Usando Comandos Go CLI
Caso não tenha o `make` instalado:

*   **Testes:** `go test -v ./...`
*   **Interativo:** `go run ./cmd/minicompiler -i`

---

## 🛠️ Desenvolvimento e Alterações

### Alterando a Gramática
Se você modificar o arquivo `pkg/compiler/parser.y`, precisará gerar o código do parser novamente:

1.  Instale o `goyacc`:
    ```bash
    go install golang.org/x/tools/cmd/goyacc@latest
    ```
2.  Gere o código:
    ```bash
    make generate
    ```

### Exemplos de Código Aceitos
```text
int x = 10;
total = (a + b) * 5;
int resultado = x + y + z;
```

---

## 📝 Pipeline de Execução

1.  **Código Fonte** (Texto)
2.  **Lexer** (`pkg/compiler/lexer.go`): Converte texto em Tokens.
3.  **Parser** (`pkg/compiler/parser.go`): Valida a estrutura gramatical.
4.  **Resultado**: Sucesso (Programa Aceito) ou Erro Sintático.
