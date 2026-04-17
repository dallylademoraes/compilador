# Minicompiler – Disciplina de Compiladores
**Gramáticas, Árvores e Parsing**

---

## Como executar no Windows

Siga estes passos para rodar o projeto no seu terminal (PowerShell ou VS Code Terminal).

### 1. Pré-requisitos
Certifique-se de ter o **Go** instalado. No Windows, o jeito mais rápido via terminal é:
```powershell
winget install GoLang.Go
```
*Após instalar, feche e abra o VS Code novamente para que o comando `go` seja reconhecido.*

### 2. Gerar o Executável
No Windows, precisamos da extensão `.exe` para que o sistema identifique o arquivo corretamente:
```powershell
go build -o minicompiler.exe .
```

### 3. Rodar a validação automática
Para testar os exemplos que já estão programados no arquivo `main.go`:
```powershell
.\minicompiler.exe
```

### 4. Modo Interativo (Testar seu próprio código)
Para digitar códigos manualmente e ver a análise léxica e sintática:
```powershell
.\minicompiler.exe -i
```

> ⚠️ **DICA PARA O WINDOWS:**
> 1. Digite seu código (ex: `x = 10;`).
> 2. Pressione **Enter**.
> 3. Pressione **`Ctrl + Z`** e depois **Enter** para finalizar a leitura e ver o resultado.
> *(O comando `Ctrl+D` mencionado em tutoriais de Linux não funciona no Windows).*

---

## Arquivos do Projeto

| Arquivo | Função |
|---|---|
| `lexer.go` | **Analisador Léxico** – Identifica tokens (palavras-chave, números, IDs). |
| `parser.y` | **Gramática** – Define as regras (o "manual de instruções" da linguagem). |
| `parser.go` | Código gerado automaticamente pelo `goyacc` a partir do `parser.y`. |
| `main.go` | **Ponto de entrada** – Coordena a execução e exibe os resultados. |

---

## Entendendo as Etapas

### 1. O Lexer (Analisador Léxico)
Identifica cada "palavra" do seu código e a classifica em categorias (Tokens).
**Exemplo:** `x = 10;`
- `x` ⮕ Identificador
- `=` ⮕ Atribuição
- `10` ⮕ Número
- `;` ⮕ Delimitador

### 2. O Parser (Analisador Sintático)
Verifica se a ordem desses tokens faz sentido segundo a gramática definida.
- `x = 10;` ✔ (Válido: segue a regra de atribuição)
- `x 10 = ;` ✗ (Erro: ordem incorreta dos tokens)

---

## Exemplos de Teste

**Código que o programa aceita:**
```text
x = 42;
int soma = 10 + 5;
total = (2 + 3) * 4;
```

**Código que gera erro:**
```text
x = 10    ← Erro: Falta o ponto-e-vírgula (;).
int = 5;  ← Erro: 'int' é reservado e não pode ser nome de variável.
x = 10 + ; ← Erro: Esperava um valor após o '+'.
```

---

## Estrutura do Pipeline

```text
Código (Texto) ──► [LEXER] ──► Tokens ──► [PARSER] ──► Resultado (Sucesso/Erro)
```

---

