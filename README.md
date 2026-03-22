# Mega Organizer (Mega Hub)

O **Mega Organizer** (ou Mega Hub) é um sistema simples e colaborativo feito em **Go** para gerenciar apostas da Mega Sena e bolões (especialmente a Mega da Virada). 
Ele permite o registro de jogos, prevenção de bilhetes repetidos (colisões), cálculo de valores de desdobramentos e visualização de estatísticas do grupo de apostadores.

## Funcionalidades

- **Registro de Jogos:** Cadastre apostas simples (6 dezenas) ou desdobramentos maiores (até 20 dezenas). O sistema possui a capacidade de preencher as dezenas restantes aleatoriamente em caso de jogos pré-fixados.
- **Prevenção de Colisão:** O sistema verifica se alguma combinação idêntica já foi apostada no banco para não haver duplicidade de jogo no bolão.
- **Dashboard e Estatísticas:** Acompanhe o total gasto, número total de jogos, e as dezenas mais "quentes" (mais apostadas).
- **Separação por Temporadas / Bolões:** Organize apostas por eventos como `Mega da Virada 2025`.
- **Histórico por Apostador:** Visualize os jogos atrelados a um *nickname* (nome de usuário).
- **Proteção de Exclusão:** Excluir apostas lançadas por engano via API precisa de um Access Token (`MONOBOLA123`).
- **Persistência Simples:** Todos os dados são salvos localmente num arquivo JSON (`apostas_db.json`).

## Tecnologias

- **Backend:** Go (Golang) com uso de pacotes nativos (`net/http`, `html/template`, `encoding/json`).
- **Frontend:** HTML, JavaScript e CSS contidos na pasta `templates/`.
- **Infraestrutura:** Docker, Docker Compose, Makefile.
- **Hot-Reload:** [Air](https://github.com/cosmtrek/air) configurado via Dockerfile para atualizar instantaneamente o servidor em ambiente de desenvolvimento quando o código é alterado.

## Rodando o Projeto

O projeto já está configurado para o Docker, permitindo que suba rapidamente sem necessariamente ter o Go instalado no host.

### Pré-requisitos
- [Docker](https://docs.docker.com/get-docker/) e Docker Compose
- GNU Make

### Comandos Make Disponíveis

Subir o projeto (executa o build e roda):
```bash
make up
```

Após subir o projeto, acesse a interface web em [http://localhost:8080](http://localhost:8080).

Outros comandos úteis:
- `make build`: Constrói a imagem Docker base (`mega-hub`).
- `make run`: Apenas roda o contêiner já formatado com o volume e mapeamento da porta `8080`.
- `make stop`: Para e remove o contêiner ativo.
- `make logs`: Verifica os logs do contêiner ativo.
- `make status`: Mostra o status do contêiner sendo executado.
- `make clean`: Cessa a execução e apaga a imagem construída localmente.
- `make reset`: **ATENÇÃO:** Para o servidor e **apaga todo o conteúdo** de `apostas_db.json`, fazendo um hard reset dos seus jogos.

### Usando puramente o Docker Compose
Se preferir não usar o Make:
```bash
docker-compose up -d --build
```

## Estrutura de Arquivos

```text
mega-organizer/
├── Dockerfile           # Imagem Go + utilitário Air instalado embutido (porta 8080)
├── Makefile             # Atalhos para lidar com a infraestrutura no Docker Compose
├── docker-compose.yml   # Volume e definição local da aplicação e mountpoint.
├── go.mod / go.sum      # Gestão de dependências do Go
├── main.go              # Servidor, roteamento, e lógica de manipulação do BD JSON
├── templates/
│   └── index.html       # Single-Page Web Frontend (Dashboard, Listagem, Envio)
└── apostas_db.json      # Arquivo gerado em runtime que funciona como Database
```

## Rotas de API

O frontend (Ajax/Fetch) se comunica com o servidor através das seguintes rotas:

- `GET /api/dados?season=X`: Pega estatísticas gerais e todas as apostas (filtrando opcionalmente por temporada).
- `GET /api/custo?qtd=N`: Calcula de antemão e devolve o preço em R$ de um jogo segundo a quantidade de dezenas.
- `GET /api/usuario/historico?nickname=Y`: Resgata todo o histórico de um dado nickname.
- `POST /api/apostar`: Rota principal para cadastrar uma nova aposta, suportando `nickname`, `qtd`, `fixos` e `season`.
- `DELETE /api/aposta/deletar?id=ID&token=MONOBOLA123`: Remove uma aposta pontualmente do banco. Requer o Token de Segurança.
