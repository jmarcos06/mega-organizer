# ======================
# Configurações
# ======================
-include .env

.PHONY: up stop logs clean reset

## Build + Run (atalho recomendado)
up:
	docker compose up -d --build

## Parar e remover containers
stop:
	docker compose down

## Logs dos containers
logs:
	docker compose logs -f

## Rodar o golangci-lint em um container isolado Docker
lint:
	@cd backend && docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run

## Rodar o linter com correção automática
lint-fix:
	@cd backend && docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run --fix

## Reiniciar
restart: stop up

## Remove tudo
clean: stop
	docker compose rm -f -v
	docker volume rm mega-organizer_mongo_data || true

## Reset total (APAGA histórico MongoDB)
reset: stop
	@echo "⚠️  Apagando volume do Banco de Dados MongoDB"
	-@docker volume rm mega-organizer_mongo_data || true
	@make up
