# ======================
# Configurações
# ======================
-include .env

.PHONY: up stop logs clean reset

## Build + Run (atalho recomendado)
up:
	@touch apostas_db.json
	docker compose up -d --build

## Parar e remover containers
stop:
	docker compose down

## Logs dos containers
logs:
	docker compose logs -f

## Reiniciar
restart: stop up

## Remove tudo
clean: stop
	docker compose rm -f -v

## Reset total (APAGA histórico local)
reset: stop
	@echo "⚠️  Apagando apostas_db.json"
	-@rm -f apostas_db.json
	@touch apostas_db.json
	@make up
