# Relatório de Implementação - (ST

## 1. Objetivo da Entrega

A entrega teve como foco evoluir o parser de uma etapa apenas validativa (aceita/rejeita entrada) para uma etapa construtiva, capaz de produzir uma Árvore Sintática Abstrata (AST) reutilizável pelas próximas fases do compilador.

A implementação também incluiu:

- integração da AST ao fluxo principal do programa como Fase 3;
- exibição visual da árvore;
- exportação da AST em DOT e PNG;
- organização dos artefatos de AST na pasta ast;
- atualização dos testes para verificar geração da AST em entradas válidas.

## 2. Resumo do Que Foi Implementado

Foram implementados os seguintes pontos:

1. Estruturas da AST no pacote compiler.
2. Ações semânticas no parser.y para construir nós durante reduções.
3. Retorno da AST pela função ExecutarParser.
4. Impressão hierárquica da AST.
5. Integração da Fase 3 no main.
6. Exportação da AST para DOT e PNG (Graphviz).
7. Organização de saída em pasta ast.
8. Atualização dos testes para validar existência da AST.

## 3. Tipo de Árvore Utilizada

Foi utilizada uma AST (Árvore Sintática Abstrata), e não uma árvore sintática concreta completa.

### 3.1 Justificativa da escolha

A AST foi escolhida porque ela representa apenas os elementos semânticos relevantes para as próximas fases (semântica, código intermediário e geração final), descartando detalhes puramente sintáticos que não agregam valor na etapa de tradução.

Exemplos de abstração adotada:

- parênteses não viram nós próprios; apenas influenciam a forma da subárvore;
- delimitadores como ponto-e-vírgula não aparecem como nó;
- operadores são representados por nós BinOp, mantendo precedência/associatividade.

Essa decisão simplifica visitantes futuros e reduz acoplamento entre parser e fases seguintes.

### 3.2 Estrutura de nós definida

Modelagem implementada:

- Node: interface base para qualquer nó.
- Stmt: interface para statements.
- Expr: interface para expressões.
- Program: raiz com lista de statements.
- Decl: declaração int x ou int x = expr.
- Assign: atribuição x = expr.
- BinOp: operação binária com operador, esquerda e direita.
- Identifier: referência a variável.
- Num: literal numérico.

Observação de decisão técnica importante:

O nó de identificador foi nomeado como Identifier (e não ID) para evitar conflito de nome com o token ID do parser gerado pelo goyacc.

## 4. Sintaxe Coberta e Mapeamento para AST

A gramática da MiniLang continuou a mesma, mas agora cada redução relevante produz um nó.

### 4.1 Regras de statements

- int id ; -> Decl com Value nulo.
- int id = expr ; -> Decl com Value apontando para expressão.
- id = expr ; -> Assign.

### 4.2 Regras de expressão

- expr + termo -> BinOp(op=+).
- expr - termo -> BinOp(op=-).
- termo * fator -> BinOp(op=*).
- termo / fator -> BinOp(op=/).
- fator -> folha ou subárvore propagada.
- ( expr ) -> retorna a própria expressão interna, sem nó de parênteses.

Com isso, a árvore final respeita precedência e associatividade já declaradas na gramática.

## 5. Ações Semânticas Implementadas no parser.y

As ações semânticas foram reescritas para construir estruturas em memória, substituindo os prints de derivação.

### 5.1 Union e tipagem de não-terminais

No bloco union foram adicionados tipos para transportar ponteiros e listas durante o parse:

- program: ponteiro para Program;
- stmt: Stmt;
- stmtList: lista de Stmt;
- expr: Expr;
- sval: string para tokens textuais.

Também foram definidos type bindings para programa, lista_stmt, stmt, expr, termo e fator.

### 5.2 Construção da raiz

Na redução de programa:

- cria Program com os statements acumulados;
- salva a referência final em YyLex.Program para recuperação após yyParse.

### 5.3 Acumulação de lista de statements

- lista_stmt -> lista_stmt stmt: append.
- lista_stmt -> stmt: inicializa slice com um elemento.

### 5.4 Construção de nós de statement

- declaração sem inicialização: Decl{Name: ...}
- declaração com inicialização: Decl{Name: ..., Value: ...}
- atribuição: Assign{Name: ..., Value: ...}

### 5.5 Construção de nós de expressão

- operações binárias constroem BinOp com Op, Left e Right;
- propagação direta em regras intermediárias evita nós redundantes;
- fator com id cria Identifier;
- fator com num cria Num;
- fator com parênteses retorna a expressão interna.
- divisão foi incluída no mesmo nível de precedência de multiplicação (%left TIMES DIV).

## 6. Integração com o Pipeline de Execução

A função ExecutarParser foi alterada para retornar dois valores:

- ponteiro para Program (AST);
- booleano de sucesso sintático.

Comportamento final:

- entrada válida: retorna AST preenchida e true;
- entrada inválida: retorna nil e false, preservando emissão de erro sintático.

No fluxo principal (main), após Fase 2:

- se parse falhar, execução interrompe com mensagem de rejeição;
- se parse passar, exibe Fase 3 e imprime a AST.
- após a impressão, o compilador cria a pasta ast e tenta exportar os arquivos ast/ast.dot e ast/ast.png.

## 7. Impressão da AST

Foi implementada uma impressão hierárquica com caracteres de árvore para facilitar inspeção visual.

Características:

- raiz Programa no topo;
- uso de galhos para filhos à esquerda/direita;
- suporte a múltiplos statements;
- saída determinística e legível para demonstração em terminal e testes.

## 8. Testes e Validação

Os testes existentes foram adaptados para o novo contrato do parser.

Validações adicionadas:

1. Para casos válidos, AST deve ser diferente de nil.
2. Para casos inválidos, AST deve ser nil.
3. Mantida a validação aceito/rejeitado já existente.
4. Incluído caso válido com divisão: int indice = custoTotal / capacidade;.

Resultado observado:

- suíte de testes passou com sucesso após regeneração do parser.

## 9. Decisões de Projeto e Racional

### 9.1 Separação entre Stmt e Expr

Separar statement e expressão por interface facilita visitantes futuros e evita estados inválidos (por exemplo, usar declaração como expressão).

### 9.2 Manter números como string no nó Num

O valor numérico foi mantido textual no nó para simplificar parser e adiar decisões de tipagem/avaliação para fases semânticas e intermediárias.

### 9.3 AST minimalista para integração com Lean e Ícaro

A árvore foi desenhada para ser simples e suficiente para:

- análise semântica (tabela de símbolos);
- geração de código intermediário por travessia de expressões e statements.

### 9.4 Compatibilidade com parser gerado

Foi adotada estratégia de armazenar o resultado final no lexer estruturado (YyLex.Program), o que mantém integração limpa com o yyParse gerado pelo goyacc.

## 10. Limitações Atuais

A implementação atual cobre a AST de forma sintática, mas ainda não executa validações semânticas (declaração prévia, redeclaração, etc.), que pertencem à etapa do Ícaro.

Também não há ainda geração de código intermediário/final, que dependem diretamente desta AST e correspondem à etapa do Lean.

Observação sobre exportação gráfica:

- ast/ast.dot é sempre gerado quando a AST existe;
- ast/ast.png depende da instalação do Graphviz e do comando dot no PATH.

## 11. Conclusão

O projeto agora possui:

- parser com construção de AST;
- visualização de árvore na Fase 3;
- contrato de retorno que desbloqueia as próximas etapas;
- testes atualizados confirmando geração correta da AST em entradas válidas.

Com isso, a base para análise semântica e geração de código está pronta para integração com os próximos módulos do pipeline.
