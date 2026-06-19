package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ResponsePayload define a estrutura estrita do JSON solicitado.
type ResponsePayload struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

// Instrumentação de Métricas (Padrão Prometheus)
var (
	// httpRequestsTotal atende ao requisito de monitoramento de "volume de requisições".
	httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Volume total de requisicoes HTTP no endpoint /projeto-korp",
	})

	// serviceAvailability atende ao requisito de monitoramento de "disponibilidade do serviço".
	serviceAvailability = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "service_availability_status",
		Help: "Status de disponibilidade do servico (1 = UP, 0 = DOWN)",
	})
)

func projetoKorpHandler(w http.ResponseWriter, r *http.Request) {
	// Incrementa a métrica de volume a cada chamada
	httpRequestsTotal.Inc()

	// O desafio especifica o verbo GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Resolução dinâmica do horário atual em UTC.
	horarioAtual := time.Now().UTC().Format(time.RFC3339)

	payload := ResponsePayload{
		Nome:    "Projeto Korp",
		Horario: horarioAtual,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func main() {
	// Inicializa a métrica de disponibilidade sinalizando que a aplicação está operante
	serviceAvailability.Set(1)

	// Definição das rotas
	http.HandleFunc("/projeto-korp", projetoKorpHandler)
	http.Handle("/metrics", promhttp.Handler())

	porta := ":8080"
	log.Printf("Servidor HTTP inicializado e escutando na porta %s", porta)

	// Inicia o servidor e altera o status de disponibilidade em caso de falha crítica
	if err := http.ListenAndServe(porta, nil); err != nil {
		serviceAvailability.Set(0)
		log.Fatalf("Falha crítica no servidor: %v", err)
	}
}
