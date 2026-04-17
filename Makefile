.PHONY: build run test clean generate

# Gera o parser usando goyacc
generate:
	goyacc -o pkg/compiler/parser.go -p yy pkg/compiler/parser.y

# Compila o projeto
build: generate
	go build -o minicompiler ./cmd/minicompiler

# Roda o modo interativo
run: build
	./minicompiler -i

# Roda os testes de validação (exemplos)
test:
	go test -v ./...

# Limpa o executável e arquivos temporários
clean:
	rm -f minicompiler pkg/compiler/y.output pkg/compiler/parser.go
