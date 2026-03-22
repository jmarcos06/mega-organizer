# Mega Organizer (Mega Hub)

Um organizador simples e ágil focado na Mega Sena e Mega da Virada para registrar e controlar as apostas do bolão da firma ou de amigos.

## Nova Arquitetura

O projeto foi inteiramente refatorado e dividido em duas camadas principais:
- **Backend (Go com Arquitetura DDD)**: Responsável pelas regras de negócio como a validação de dezenas, cálculos de custo da aposta, e verificação de colisões (duplicidade) de bilhetes.
- **Frontend (Vue.js + Vite)**: Aplicação Web SPA (Single Page Application) de resposta rápida. Consome a API REST do backend usando funcionalidades nativas de reatividade, sendo totalmente estilizado com as classes do Tailwind CSS.

## Funcionalidades
- **Registro de Apostas:** Aceita jogos simples (6) e desdobramentos (até 20 dezenas) e seleciona as dezenas remanescentes de forma aleatória em bilhetes com dezenas fixas.
- **Cálculo Automatizado e Prevenção de Colisão:** Calcula os preços e alerta automaticamente caso a combinação exata de botões já tenha sido registrada por outro usuário do site.
- **Filtros e Histórico:** Separação de jogos por temporadas/eventos (Ex: *Mega da Virada 2025*), visualização de dezenas mais "quentes", e histórico clicável para inspecionar jogos de um Nickname específico.
- **Exclusão Transacional:** Jogos excluidos contam com um JWT falso (`MONOBOLA123`) para proteger contra uso malicioso de bots na API Pública.

## Como Rodar o Projeto

Usamos `docker compose` para simplificar a inicialização simultânea de todas as aplicações dependentes. Um container do **MongoDB** é inicializado junto com a API e o Frontend, simulando idênticamente o ambiente de produção Cloud.

### Requisitos
- Docker (com Compose V2 habilitado)
- GNU Make

### Subindo a Aplicação
Através do terminal na pasta do projeto, execute:
```bash
make up
```
Isso acionará o download das dependências no node e go e subirá as duas engrenagens. 

Acesse a **interface do sistema** em: [http://localhost:5173](http://localhost:5173).

### Outros Comandos úteis:
- `make stop`: Desliga todos os containeres do projeto.
- `make logs`: Verifica os logs e outputs (Vite/Go).
- `make lint`: Roda o `golangci-lint` oficial via Docker no código do backend para checar a qualidade da aplicação.
- `make lint-fix`: Roda o linter corrigindo automaticamente infrações simples de formatação go.
- `make clean`: Limpa os containers criados e mata os volumes pendentes.
- `make reset`: **CUIDADO.** Elimina permanentemente o volume de persistência offline do MongoDB e inicializa a base de dados zerada na próxima execução.

## Estrutura de Pastas e Códigos

A configuração prioriza um fluxo unificado de dev onde você altera os módulos dinamicamente.

```text
mega-organizer/
├── backend/                  # API Rest (Golang + Arquitetura DDD)
│   ├── cmd/server/           # Ponto de entrada (Main) da aplicação
│   ├── internal/
│   │   ├── application/      # Casos de uso (CreateBet, GetUserHistory...)
│   │   ├── domain/           # Entidades (Bet, Stats...) e lógica matemática principal
│   │   ├── infrastructure/   # Repositórios (Integração e consultas com o MongoDB)
│   │   └── interfaces/http/  # Roteamento da API, Middlewares de Auth e Cors
│   ├── Dockerfile            # Container isolado acoplado com o `Air` (Hot-Reload para o Go)
│   └── .air.toml             # Config. do binário temporário de execução
├── frontend/                 # Client consumid0r construído via Vite
│   ├── src/
│   │   ├── main.js           # Inicialização principal do Node/Vue
│   │   ├── App.vue           # Componente Raiz da Aplicação (Layout + Requests)
│   │   └── style.css         # Assets base do Vite
│   ├── index.html            # Montagem estrutural que injeta as cdn do Tailwind e ícones
│   ├── vite.config.js        # Configuração e direcionamento do proxy da `/api` pra porta `8080`
│   └── Dockerfile            # Container para a exposição do Vite Client
├── docker-compose.yml        # Acoplamento de Rede e Banco (MongoDB local)
├── Makefile                  # Comandos de orquestração e Linting
└── backend/.golangci.yml     # Configurações do linter GolangCI
```

## Referência da Roteamento e Mapeamento Back -> Front

Caso use insomnia/postman, o backend escuta abertamente em `:8080/api`:
- `GET /api/dados?season=...`: Retorna todas as métricas dos jogos registrados.
- `GET /api/custo?qtd=6`: Utilizado pra prever o custo total de x números apostados sob o preço da aposta simples.
- `GET /api/usuario/historico?nickname=XYZ`: Usado na modal do frontend.
- `POST /api/apostar`: Recebe payload JSON gerando dezenas complementares aos preenchimentos em Array. Efetua a checagem no banco.
- `DELETE /api/aposta/deletar`: Elimina um array do banco. Requer o parâmetro Header string.
