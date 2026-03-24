.PHONY: run build tidy clean

run:
	go run ./cmd/api

build:
	go build -o tech-memo ./cmd/api

tidy:
	go mod tidy

clean:
	rm -f tech-memo tech-memo.exe tech_memo.db
