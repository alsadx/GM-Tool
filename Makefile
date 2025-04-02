.PHONY:
.SILENT:

run:
	go build -o gm ./cmd/app/main.go
	./gm