FROM golang:1.25-alpine

WORKDIR /app

# Dependências mínimas
RUN apk add --no-cache git

# Instala o Air
RUN go install github.com/cosmtrek/air@v1.52.0

# Garante que o Air esteja no PATH
ENV PATH="/go/bin:${PATH}"

# Copia apenas módulos primeiro (para caching eficiente)
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código
COPY . .

# Exponha a porta
EXPOSE 8080

# Air roda direto
CMD ["air"]
