# 📝 Relatório Técnico: MiniCompilador (Foco: Ícaro — Tabela de Símbolos & Análise Semântica)

Este relatório descreve tecnicamente a implementação da **Fase 4 (Análise Semântica)** no MiniCompilador, com ênfase na **Tabela de Símbolos** e no **Visitor Semântico**.

---

## 1. Objetivo da Entrega

A entrega teve como objetivo introduzir validação semântica sobre a AST já construída pelas fases anteriores, para impedir aceitação de programas sintaticamente válidos, porém semanticamente inválidos.

Escopo atendido:

1. Implementar estrutura de **Tabela de Símbolos**.
2. Implementar um **analisador semântico** baseado em Visitor.
3. Integrar a Fase 4 ao fluxo do executável principal.
4. Cobrir os cenários semânticos básicos com testes automatizados.

---

## 2. Arquitetura da Solução

### 2.1 Componentes adicionados

- `pkg/compiler/simbolos.go`
- `pkg/compiler/semantico.go`

### 2.2 Posição no pipeline

A análise semântica ocorre após parser/AST:

1. Léxico (`lexer.go`)
2. Sintático + AST (`parser.y` / `parser.go` / `ast.go`)
3. **Semântico (`semantico.go` + `simbolos.go`)**
4. Código intermediário e código alvo (fases posteriores)

No `cmd/minicompiler/main.go`, a Fase 4 é executada via:

```go
tabela, erros := compiler.AnalisarSemantica(ast)
```

Se houver erros, o pipeline é interrompido antes das etapas seguintes.

---

## 3. Tabela de Símbolos (`simbolos.go`)

### 3.1 Estruturas de dados

```go
type Simbolo struct {
    Nome  string
    Tipo  string
    Valor string
}

type TabelaSimbolos struct {
    tabela map[string]Simbolo
}
```

### 3.2 Operações implementadas

1. `NovaTabelaSimbolos()`  
   Inicializa a estrutura de armazenamento.

2. `Inserir(nome, tipo string)` / `InserirComValor(nome, tipo, valor string)`  
   Registra símbolo novo; falha em redeclaração.

3. `Atribuir(nome, valor string)`  
   Atualiza símbolo existente; falha se variável não declarada.

4. `Buscar(nome string) (Simbolo, bool)`  
   Consulta de existência e dados do símbolo.

5. `Imprimir()` / `String()`  
   Geração de tabela textual formatada para saída de terminal.

### 3.3 Contrato de erro

Mensagens padronizadas:

- `[Erro semântico] Variável '<id>' já foi declarada`
- `[Erro semântico] Variável '<id>' usada antes de ser declarada`

---

## 4. Analisador Semântico (`semantico.go`)

### 4.1 Estrutura do analisador

```go
type AnalisadorSemantico struct {
    tabela *TabelaSimbolos
    erros  []string
}
```

A função pública de entrada é:

```go
func AnalisarSemantica(raiz *Program) (*TabelaSimbolos, []string)
```

Ela instancia o analisador, executa `Check` e retorna tabela + lista de erros.

### 4.2 Estratégia de travessia

A travessia segue o padrão Visitor já definido na AST:

- `VisitProgram`: percorre todos os statements.
- `VisitDecl`: valida expressão inicial e registra símbolo.
- `VisitAssign`: valida expressão e atualiza símbolo existente.
- `VisitBinOp`: visita recursivamente esquerda e direita.
- `VisitIdentifier`: valida uso de variável previamente declarada.
- `VisitNum`: sem validações adicionais.

### 4.3 Regra de acúmulo de erros

Os erros são acumulados em `[]string`, sem abortar no primeiro problema.  
Isso permite retornar um conjunto de inconsistências semânticas encontradas na análise.

### 4.4 Representação de valor na tabela

A função `formatarExpr` serializa expressões para texto (ex.: `(y + x)`), que é o formato armazenado em `Simbolo.Valor`.  
Não há avaliação numérica de expressões nesta fase.

---

## 5. Regras Semânticas Cobertas

No estado atual da implementação, a Fase 4 cobre:

1. **Uso de variável antes da declaração**.
2. **Redeclaração de variável**.
3. **Atribuição para variável não declarada**.

Essas regras atendem o núcleo funcional da tabela de símbolos para integração com as fases seguintes.

---

## 6. Integração com o Executável Principal

No fluxo de execução (`cmd/minicompiler/main.go`):

1. AST é produzida na Fase 2/3.
2. Fase 4 chama `AnalisarSemantica`.
3. Se houver erros:
   - imprime mensagens semânticas;
   - encerra com rejeição semântica.
4. Se não houver erros:
   - imprime confirmação;
   - exibe tabela de símbolos formatada.

Esse comportamento estabelece um gate semântico antes da geração de código.

---

## 7. Validação por Testes

Arquivo principal: `cmd/minicompiler/main_test.go`  
Teste dedicado: `TestAnaliseSemantica`.

Cenários validados:

1. Programa semanticamente válido.
2. Uso antes da declaração.
3. Atribuição sem declaração.
4. Redeclaração de variável.

Asserções realizadas:

- correspondência entre sucesso/falha semântica esperada e obtida;
- presença de mensagens de erro esperadas;
- presença de símbolos esperados na tabela para casos válidos.

---

## 8. Impacto Técnico da Entrega

A contribuição do Ícaro transformou o projeto de um pipeline apenas léxico/sintático em um compilador com validação semântica mínima operacional.

Benefícios diretos:

1. detecção antecipada de inconsistências de identificação de variáveis;
2. construção de ambiente simbólico reutilizável por fases posteriores;
3. contrato claro de erro semântico para CLI e testes.

---

## 9. Limitações Atuais e Evolução Natural

Limitações observadas no estado atual:

- sem inferência/verificação de tipos além do tipo `int` implícito;
- sem avaliação constante para regras semânticas avançadas;
- sem escopos aninhados (modelo atual é tabela global única).

Próximas evoluções naturais:

1. checagens semânticas adicionais (ex.: regras sobre operações específicas);
2. suporte a escopo léxico/blocos;
3. enriquecimento da tabela com metadados para otimização e geração de código.

---

*Relatório técnico consolidado a partir da implementação atual de `simbolos.go`, `semantico.go` e integração no `main.go`.*
