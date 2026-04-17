%{
package compiler

import "fmt"
%}

%union {
    sval string
}

%token <sval> KEYWORD_INT
%token <sval> ID
%token <sval> ASSIGN
%token <sval> PLUS
%token <sval> MINUS
%token <sval> TIMES
%token <sval> SEMI
%token <sval> LPAREN
%token <sval> RPAREN
%token <sval> NUM
%token <sval> ERROR

%left PLUS MINUS
%left TIMES

%%

// Gramática da minilinguagem:
//
// programa   → lista_stmt
// lista_stmt → lista_stmt stmt | stmt
// stmt       → int id ;
//            | int id = expr ;
//            | id = expr ;
// expr       → expr + termo | expr - termo | termo
// termo      → termo * fator | fator
// fator      → id | num | ( expr )

programa
    : lista_stmt
        { fmt.Println("\n✔  Análise sintática concluída: programa aceito.") }
    ;

lista_stmt
    : lista_stmt stmt
    | stmt
    ;

stmt
    : KEYWORD_INT ID SEMI
        {
            fmt.Printf("[Parser] Declaração simples → int %s\n", $2)
        }
    | KEYWORD_INT ID ASSIGN expr SEMI
        {
            fmt.Printf("[Parser] Declaração com atribuição → int %s = <expr>\n", $2)
        }
    | ID ASSIGN expr SEMI
        {
            fmt.Printf("[Parser] Atribuição → %s = <expr>\n", $1)
        }
    ;

expr
    : expr PLUS termo
        { fmt.Println("[Parser] Redução: expr → expr + termo") }
    | expr MINUS termo
        { fmt.Println("[Parser] Redução: expr → expr - termo") }
    | termo
    ;

termo
    : termo TIMES fator
        { fmt.Println("[Parser] Redução: termo → termo * fator") }
    | fator
    ;

fator
    : ID
        { fmt.Printf("[Parser] Fator → id(%s)\n", $1) }
    | NUM
        { fmt.Printf("[Parser] Fator → num(%s)\n", $1) }
    | LPAREN expr RPAREN
        { fmt.Println("[Parser] Fator → ( expr )") }
    ;

%%
