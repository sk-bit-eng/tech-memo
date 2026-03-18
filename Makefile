.PHONY: run build tidy clean

run:
	go run ./cmd/api

build:
	go build -o tech-memo ./cmd/api

tidy:
	go mod tidy

clean:
	powershell -Command "if (Test-Path tech-memo.exe) { Remove-Item tech-memo.exe }; if (Test-Path tech-memo) { Remove-Item tech-memo }; if (Test-Path tech_memo.db) { Remove-Item tech_memo.db }"
