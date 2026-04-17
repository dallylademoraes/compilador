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
