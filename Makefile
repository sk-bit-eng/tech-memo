.PHONY: run build tidy clean

run:
	go run ./cmd/main.go

build:
	go build -o tech-memo ./cmd/main.go

tidy:
	go mod tidy

clean:
	rm -f tech-memo tech_memo.db
