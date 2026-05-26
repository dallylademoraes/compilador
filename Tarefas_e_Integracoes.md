## 👥 Divisão de Tarefas

---

### 👤 João Walcacer — AST (Árvore Sintática Abstrata)
> Seções **4 e 5** do trabalho · Análise Sintática + Tradução Dirigida por Sintaxe

**Contexto:** O parser atual só imprime mensagens no console. Precisa evoluir para construir uma árvore de nós em memória que represente o programa.

**Tarefas:**

1. Criar `ast.go` com as structs de nó:
   - `NóPrograma` — raiz da árvore, contém lista de statements
   - `NóDecl` — declaração `int x` ou `int x = expr`
   - `NóAtrib` — atribuição `x = expr`
   - `NóBinOp` — operação binária com campos `Op`, `Esquerda`, `Direita`
   - `NóID` — referência a variável (folha)
   - `NóNum` — número literal (folha)

2. Modificar as **ações semânticas do `parser.y`** para popular a AST em vez de só imprimir.

3. Implementar `ImprimirAST(nó Nó, indent int)` com saída visual:
   ```
   Programa
   └── Decl: int soma
       └── BinOp: +
           ├── ID: a
           └── ID: b
   ```

4. Adicionar a exibição da AST como **Fase 3** no `main.go` e nos testes.

**Entrega:** `ExecutarParser` deve retornar a AST populada para o Lean usar.

> ⚠️ **Prioridade máxima** — o Lean está bloqueado até isso estar pronto.

---

### 👤 Ícaro — Tabela de Símbolos + Análise Semântica
> Seções **5 e 7** do trabalho · Tradução Dirigida por Sintaxe + Ambientes de Execução

**Contexto:** A linguagem tem variáveis tipadas (`int`). Precisa rastrear quais foram declaradas, seus tipos e valores para detectar erros semânticos.

**Tarefas:**

1. Criar `simbolos.go` com a struct `TabelaSimbolos`:
   ```go
   type Simbolo struct {
       Nome  string
       Tipo  string // "int"
       Valor string // "" se não inicializado
   }
   type TabelaSimbolos struct {
       tabela map[string]Simbolo
   }
   ```

2. Implementar os métodos:
   - `Inserir(nome, tipo string)` — registra variável; erro se já declarada
   - `Atribuir(nome, valor string)` — atualiza valor; erro se não declarada
   - `Buscar(nome string) (Simbolo, bool)` — consulta variável
   - `Imprimir()` — exibe a tabela formatada no terminal

3. Criar `semantico.go` com visitor que percorre a AST do João:
   - `NóDecl` → chama `Inserir`
   - `NóAtrib` → verifica se a variável existe
   - `NóID` em expressão → verifica se foi declarado

4. Reportar erros semânticos:
   ```
   [Erro semântico] Variável 'x' usada antes de ser declarada
   [Erro semântico] Variável 'y' já foi declarada
   ```

5. Adicionar exibição da tabela como **Fase 4** no `main.go`.

**Entrega:** tabela de símbolos funcionando para o Lean usar na geração de código.

> ⚠️ **Prioridade máxima** — o Lean está bloqueado até isso estar pronto.

---

### 👤 Lean — Código Intermediário + Geração de Código Final
> Seções **6 e 8** do trabalho · Geração de Código Intermediário + Geração de Código

**Contexto:** Com a AST (João) e a tabela de símbolos (Ícaro) prontas, transforma o programa em código executável em dois passos.

**Tarefas — Parte 1: Código Intermediário (`intermediario.go`)**

1. Gerador de **código de três endereços** percorrendo a AST:
   - Expressão binária → temporário: `t1 = a + b`
   - Declaração → `int soma` + `soma = t1`
   - Atribuição simples → `x = 42`

2. Exemplo para `int soma = a + b * c;`:
   ```
   t1 = b * c
   t2 = a + t1
   soma = t2
   ```

3. Structs de instrução:
   ```go
   type Instrucao struct {
       Resultado string
       Op        string
       Arg1      string
       Arg2      string
   }
   ```

4. Exibir como **Fase 5** no `main.go`.

**Tarefas — Parte 2: Geração de Código Final (`gerador.go`)**

1. Traduzir instruções para **código Python válido**:
   - `int x;` → `x = None`
   - `int soma = a + b;` → `soma = a + b`
   - `x = 42;` → `x = 42`

2. Salvar com flag `-o`:
   ```bash
   go run . -i -o saida.py
   ```

3. Exibir como **Fase 6** no `main.go`.

**Entrega:** lista de `Instrucao` disponível para a Dallyla otimizar.

---

### 👤 Dallyla — Otimizações + Relatório + Apresentação
> Seções **9, 10, 11 e 12** do trabalho

**Contexto:** Otimizações aplicadas sobre as instruções de três endereços do Lean. O relatório cobre todas as 12 seções.

**Tarefas — Implementação (`otimizador.go`)**

1. **Propagação de constantes:** se `t1 = 5 + 3`, substituir diretamente por `t1 = 8`.

2. **Eliminação de código morto:** se um temporário `t1` nunca é usado, removê-lo da lista.

3. Exibir o código intermediário **antes e depois** da otimização.

**Tarefas — Relatório**

| Seção | O que escrever |
|---|---|
| 1 – Introdução | Descrição da MiniLang, exemplos de programas válidos |
| 2 – Gramática e Autômatos | BNF do `parser.y` + **diagrama do autômato finito léxico** (ver tarefa abaixo ⚠️) |
| 3 – Análise Léxica | Explicar o `lexer.go`, mostrar tabela de tokens de exemplo |
| 4 – Análise Sintática | Explicar o parser LALR(1), mostrar exemplo de AST gerada |
| 5 – Tradução Dirigida | Explicar as ações semânticas e a tabela de símbolos |
| 6 – Código Intermediário | Mostrar exemplos de três endereços gerados |
| 7 – Ambiente de Execução | Descrever como a tabela de símbolos gerencia memória |
| 8 – Geração de Código | Mostrar exemplo de código Python gerado |
| 9 – Otimizações | Mostrar antes/depois da propagação de constantes e eliminação de mortos |
| 10 – Paralelismo | **Grafo de dependências entre instruções** (ver tarefa abaixo ⚠️) |
| 11 – Localidade | Discutir como reordenar instruções poderia melhorar cache |
| 12 – Aplicações | Relacionar com JIT (Java/.NET), compiladores para GPU, sistemas embarcados |

> ⚠️ **Tarefa obrigatória — Seção 2: Autômato Finito Léxico**  
> O professor pede explicitamente a **representação do autômato finito** usado na análise léxica, não só a gramática. Precisa ser um diagrama de estados mostrando as transições para cada categoria de token (`int`, identificadores, números, operadores, delimitadores). Pode ser desenhado à mão, gerado com uma ferramenta (ex: draw.io, Graphviz) ou incluído no relatório como imagem. **Só mencionar no texto não é suficiente.**

> ⚠️ **Tarefa obrigatória — Seção 10: Grafo de Dependências**  
> O professor pede uma **representação que evidencie dependências entre instruções**, não só análise teórica. Para o código intermediário gerado pelo Lean, Dallyla deve produzir um grafo simples onde cada nó é uma instrução e as arestas indicam qual instrução depende do resultado de outra. Exemplo para `t1 = b * c / t2 = a + t1 / soma = t2`:
> ```
> [t1 = b * c] ──→ [t2 = a + t1] ──→ [soma = t2]
> ```
> Pode ser uma imagem, um diagrama ASCII no relatório, ou gerado por código. **Instruções sem dependência entre si ficam em paralelo no grafo.**

**Tarefas — Apresentação**

- Slide 1: Título e integrantes
- Slide 2: Visão geral da arquitetura
- Slides 3–8: Uma fase por slide (léxico, sintático, AST, tabela, intermediário, código final)
- Slide 9: Otimizações com exemplo antes/depois
- Slide 10: Demo ao vivo + conclusão

> 💡 As seções 1–5 do relatório já podem ser escritas agora com o que existe no projeto.

---

## 🔗 Pontos de integração

```
João (AST) ──────────────────────────────────→ Lean (usa a AST pra gerar código)
                                             ↗
Ícaro (Tabela de Símbolos) ──────────────────

Lean (Código Intermediário) ─────────────────→ Dallyla (aplica otimizações)
```

| Handoff | Quem entrega | Quem recebe | Pré-requisito |
|---|---|---|---|
| `NóPrograma` populado + `ImprimirAST` | João | Lean | Structs de nó definidas, `parser.y` gerando a árvore |
| `TabelaSimbolos` + visitor semântico | Ícaro | Lean | Tabela funcionando, erros semânticos reportados |
| Lista de `Instrucao` (três endereços) | Lean | Dallyla | Gerador percorrendo a AST corretamente |