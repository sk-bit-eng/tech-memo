.PHONY: run build tidy clean docker-up docker-down docker-seed docker-test

run:
	go run ./cmd/api

build:
	go build -o tech-memo ./cmd/api

tidy:
	go mod tidy

clean:
	rm -f tech-memo tech-memo.exe tech_memo.db

docker-up:
	docker compose up -d sqlserver app

docker-down:
	docker compose down -v

docker-seed:
	docker compose --profile tools run --rm seed

docker-test:
	docker compose --profile test run --rm test
