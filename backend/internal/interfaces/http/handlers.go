package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	appBet "mega-play/internal/application/bet"
	domainBet "mega-play/internal/domain/bet"
)

const AccessToken = "MONOBOLA123"

func CorsMiddlewareDynamic(next http.HandlerFunc, allowedOrigins []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allow := false
		if len(allowedOrigins) == 0 {
			allow = true // Local dev
		} else {
			for _, o := range allowedOrigins {
				if o == origin || strings.HasSuffix(origin, ".vercel.app") {
					allow = true
					break
				}
			}
		}

		if allow {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "https://monobola.com")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token != AccessToken {
			http.Error(w, "Acesso não autorizado", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func HandleGetData(useCases *appBet.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentSeason := r.URL.Query().Get("season")
		stats, err := useCases.GetDashboardStats(currentSeason)
		if err != nil {
			http.Error(w, `{"error": "Erro interno"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

func HandleGetCost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qtdStr := r.URL.Query().Get("qtd")
		qtd, err := strconv.Atoi(qtdStr)
		if err != nil || qtd < 6 || qtd > 20 {
			http.Error(w, `{"error": "Quantidade inválida"}`, http.StatusBadRequest)
			return
		}

		custo := domainBet.CalcularCusto(qtd)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{"custo": custo})
	}
}

func HandleGetUserHistory(useCases *appBet.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.URL.Query().Get("nickname")
		if nickname == "" {
			http.Error(w, `{"error": "Nickname é obrigatório"}`, http.StatusBadRequest)
			return
		}

		hist, err := useCases.GetUserHistory(nickname)
		if err != nil {
			http.Error(w, `{"error": "Erro interno"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hist)
	}
}

func HandleDeleteBet(useCases *appBet.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, `{"error": "ID é obrigatório"}`, http.StatusBadRequest)
			return
		}

		if err := useCases.DeleteBet(id); err != nil {
			http.Error(w, `{"error": "Aposta não encontrada"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func HandlePostBet(useCases *appBet.UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}

		var input struct {
			Nickname string `json:"nickname"`
			QtdStr   string `json:"qtd"`
			FixosStr string `json:"fixos"`
			Season   string `json:"season"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, `{"error": "JSON inválido"}`, http.StatusBadRequest)
			return
		}

		qtd, _ := strconv.Atoi(input.QtdStr)
		var fixos []int
		if input.FixosStr != "" {
			input.FixosStr = strings.ReplaceAll(input.FixosStr, ",", " ")
			parts := strings.Fields(input.FixosStr)
			for _, p := range parts {
				val, err := strconv.Atoi(strings.TrimSpace(p))
				if err == nil && val >= 1 && val <= 60 {
					fixos = append(fixos, val)
				}
			}
		}

		bet, colisao, quem, err := useCases.CreateBet(input.Nickname, qtd, fixos, input.Season)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"aposta":      bet,
			"colisao":     colisao,
			"colisao_com": quem,
		})
	}
}
