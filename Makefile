build:
	go build -o bin/repl ./cmd/repl

repl: build
	./bin/repl
