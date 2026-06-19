# ==========================================
# Stage 1: Build
# ==========================================
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Cópia do manifesto do módulo (go.sum gerado automaticamente pelo go mod download)
COPY app/go.mod ./
RUN go mod download

# Cópia do código-fonte do subdiretório
COPY app/ .

# Compilação estática do binário (sem dependências dinâmicas CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-server-projeto-korp .

# ==========================================
# Stage 2: Runtime
# ==========================================
FROM alpine:latest

# Instalação de certificados de autoridade e base de fuso horário (tzdata)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Cópia do binário isolado gerado no Stage 1
COPY --from=builder /build/http-server-projeto-korp .

EXPOSE 8080

CMD ["./http-server-projeto-korp"]
