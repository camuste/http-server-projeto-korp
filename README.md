# http-server-projeto-korp

Serviço HTTP em Go com monitoramento via Prometheus e Grafana, proxy reverso com NGINX e automação com Ansible.

---

## Arquitetura

```
Usuário
  │
  │ porta 80
  ▼
NGINX (proxy reverso)
  │
  │ porta 8080 (rede interna Docker)
  ▼
http-server-projeto-korp (Go)
  │
  │ /metrics
  ▼
Prometheus (coleta métricas a cada 15s)
  │
  ▼
Grafana (visualização dos dados)
```

Todos os containers compartilham a rede bridge `korp-net`. A porta 8080 da aplicação Go nunca é exposta ao host — o NGINX é o único ponto de entrada público.

---

## Estrutura de pastas

```
http-server-projeto-korp/
├── app/
│   ├── go.mod               # dependências Go
│   └── main.go              # servidor HTTP + métricas Prometheus
├── nginx/
│   └── http-server-projeto-korp.conf  # configuração do proxy reverso
├── prometheus/
│   └── prometheus.yml       # configuração de coleta de métricas
├── grafana/
│   ├── provisioning/
│   │   ├── datasources/
│   │   │   └── datasources.yml  # conecta Grafana ao Prometheus
│   │   └── dashboards/
│   │       └── dashboards.yml   # aponta para os arquivos de dashboard
│   └── dashboards/
│       └── http-server-projeto-korp-dashboard.json  # dashboard pronto
├── ansible/
│   ├── inventory.ini        # alvo de execução do Ansible
│   ├── requirements.yml     # coleções Ansible necessárias
│   └── playbook.yml         # automação completa do ambiente
├── Dockerfile               # build multi-stage da aplicação Go
└── docker-compose.yml       # orquestração de todos os containers
```

---

## Pré-requisitos

### Para executar com Docker Compose
- [Docker](https://docs.docker.com/get-docker/) instalado e rodando
- Docker Compose (incluído no Docker Desktop)

### Para executar com Ansible
- Sistema Linux (Ubuntu/Debian)
- Python 3 instalado
- Ansible instalado: `pip install ansible`

---

## Como executar com Docker Compose

```bash
# 1. Clone o repositório
git clone https://github.com/camuste/http-server-projeto-korp.git
cd http-server-projeto-korp

# 2. Suba o ambiente completo
docker compose up --build -d

# 3. Teste o serviço
curl http://localhost:80/projeto-korp
```

Resposta esperada:
```json
{"nome":"Projeto Korp","horario":"2026-06-18T22:00:00Z"}
```

Para encerrar:
```bash
docker compose down
```

---

## Como executar com Ansible

```bash
# 1. Instale a coleção necessária
ansible-galaxy collection install -r ansible/requirements.yml

# 2. Execute o playbook (provisiona tudo com um único comando)
ansible-playbook -i ansible/inventory.ini ansible/playbook.yml
```

O playbook irá:
1. Instalar o Docker
2. Criar a rede Docker
3. Copiar os arquivos do projeto
4. Fazer o build e subir os containers
5. Aguardar o serviço ficar disponível
6. Validar com uma requisição HTTP e exibir a resposta no console

---

## Como testar

### Endpoint principal
```bash
curl http://localhost:80/projeto-korp
```

### Métricas expostas
A porta 8080 não é publicada no host (de propósito — só o NGINX é exposto), então o `/metrics` é acessado pela mesma porta 80, via proxy:
```bash
curl http://localhost/metrics
```

### Prometheus
Acesse `http://localhost:9090` → Status → Targets

O target `http-server-projeto-korp:8080` deve aparecer com status **UP**.

### Grafana
Acesse `http://localhost:3000`

- Usuário: `admin`
- Senha: `admin`

O dashboard **"Projeto Korp — HTTP Server"** é carregado automaticamente com 3 painéis:
- **Status do Serviço** — UP (verde) ou DOWN (vermelho)
- **Total de Requisições** — contador acumulado
- **Taxa de Requisições/segundo** — gráfico em tempo real

---

## Endpoints da aplicação

| Endpoint | Método | Descrição |
|---|---|---|
| `/projeto-korp` | GET | Retorna JSON com nome e horário UTC |
| `/metrics` | GET | Métricas no formato Prometheus |

---

## Métricas implementadas

| Métrica | Tipo | Descrição |
|---|---|---|
| `http_requests_total` | Counter | Total de requisições recebidas |
| `service_availability_status` | Gauge | 1 = UP, 0 = DOWN |

---

## Decisões técnicas

**Por que multi-stage build no Dockerfile?**
O Stage 1 usa a imagem `golang` (~300 MB) apenas para compilar. O Stage 2 usa `alpine` (~5 MB) apenas para executar o binário. A imagem final fica com ~15 MB.

**Por que a porta 8080 não é exposta ao host?**
O desafio exige isolamento: a aplicação Go só deve ser acessível internamente, via NGINX. O `expose` no Docker Compose torna a porta visível apenas dentro da rede `korp-net`.

**Por que `promauto` em vez de `prometheus.MustRegister`?**
O `promauto` registra as métricas automaticamente na criação, eliminando a necessidade de um bloco `init()` separado — menos código com o mesmo resultado.

**Por que provisioning automático no Grafana?**
Os arquivos `datasources.yml`, `dashboards.