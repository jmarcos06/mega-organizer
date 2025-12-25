# ======================
# Configurações (env)
# ======================
-include .env

IMAGE_NAME ?= mega-hub
CONTAINER_NAME ?= mega-server
PORT ?= 8080
DATA_DIR ?= $(shell pwd)

.PHONY: build run up stop logs clean reset status

## Build da imagem
build:
	docker build -t $(IMAGE_NAME) .

## Rodar container
run:
	docker run -d \
		-p $(PORT):8080 \
		-v $(DATA_DIR):/data \
		--name $(CONTAINER_NAME) \
		$(IMAGE_NAME)

## Build + Run (atalho recomendado)
up: build run

## Parar e remover container
stop:
	-@docker stop $(CONTAINER_NAME)
	-@docker rm $(CONTAINER_NAME)

## Logs do container
logs:
	docker logs -f $(CONTAINER_NAME)

## Status do container
status:
	docker ps -a | grep $(CONTAINER_NAME) || echo "Container não existe"

## Remove container + imagem
clean: stop
	-@docker rmi $(IMAGE_NAME)

## Reset total (APAGA histórico local)
reset: stop
	@echo "⚠️  Apagando apostas_db.json"
	-@rm -f apostas_db.json
	@make up
