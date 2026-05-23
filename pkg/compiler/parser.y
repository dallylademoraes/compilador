%{
package compiler
%}

%union {
    sval     string
    program  *Program
    stmt     Stmt
    stmtList []Stmt
    expr     Expr
}

%token <sval> KEYWORD_INT
%token <sval> ID
%token <sval> ASSIGN
%token <sval> PLUS
%token <sval> MINUS
%token <sval> TIMES
%token <sval> DIV
%token <sval> SEMI
%token <sval> LPAREN
%token <sval> RPAREN
%token <sval> NUM
%token <sval> ERROR

%type <program> programa
%type <stmtList> lista_stmt
%type <stmt> stmt
%type <expr> expr termo fator

%left PLUS MINUS
%left TIMES DIV

%%

// Gramática da minilinguagem:
//
// programa   → lista_stmt
// lista_stmt → lista_stmt stmt | stmt
// stmt       → int id ;
//            | int id = expr ;
//            | id = expr ;
// expr       → expr + termo | expr - termo | termo
// termo      → termo * fator | termo / fator | fator
// fator      → id | num | ( expr )

programa
    : lista_stmt
        {
            $$ = &Program{Statements: $1}
            if l, ok := yylex.(*YyLex); ok {
                l.Program = $$
            }
        }
    ;

lista_stmt
    : lista_stmt stmt
        { $$ = append($1, $2) }
    | stmt
        { $$ = []Stmt{$1} }
    ;

stmt
    : KEYWORD_INT ID SEMI
        {
            $$ = &Decl{Name: $2}
        }
    | KEYWORD_INT ID ASSIGN expr SEMI
        {
            $$ = &Decl{Name: $2, Value: $4}
        }
    | ID ASSIGN expr SEMI
        {
            $$ = &Assign{Name: $1, Value: $3}
        }
    ;

expr
    : expr PLUS termo
        { $$ = &BinOp{Op: $2, Left: $1, Right: $3} }
    | expr MINUS termo
        { $$ = &BinOp{Op: $2, Left: $1, Right: $3} }
    | termo
        { $$ = $1 }
    ;

termo
    : termo TIMES fator
        { $$ = &BinOp{Op: $2, Left: $1, Right: $3} }
    | termo DIV fator
        { $$ = &BinOp{Op: $2, Left: $1, Right: $3} }
    | fator
        { $$ = $1 }
    ;

fator
    : ID
        { $$ = &Identifier{Name: $1} }
    | NUM
        { $$ = &Num{Value: $1} }
    | LPAREN expr RPAREN
        { $$ = $2 }
    ;

%%
