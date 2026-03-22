package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	appBet "mega-play/internal/application/bet"
	"mega-play/internal/infrastructure/storage"
	interfacesHttp "mega-play/internal/interfaces/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mongoURI := os.Getenv("DATABASE_URL")
	if mongoURI == "" {
		mongoURI = "mongodb://database:27017"
	}

	// Initialize MongoDB Repository
	repo, err := storage.NewMongoRepository(mongoURI)
	if err != nil {
		log.Fatalf("Falha ao conectar no MongoDB: %v", err)
	}

	if err := repo.Load(); err != nil {
		log.Printf("⚠️ Aviso: Falha no Ping do MongoDB: %v", err)
	} else {
		bets, _ := repo.GetAll()
		log.Printf("📂 %d apostas carregadas do histórico da cloud.", len(bets))
	}

	// Initialize UseCases
	useCases := appBet.NewUseCases(repo)

	// Register Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/dados", interfacesHttp.HandleGetData(useCases))
	mux.HandleFunc("/api/apostar", interfacesHttp.HandlePostBet(useCases))
	mux.HandleFunc("/api/custo", interfacesHttp.HandleGetCost())
	mux.HandleFunc("/api/usuario/historico", interfacesHttp.HandleGetUserHistory(useCases))
	mux.HandleFunc("/api/aposta/deletar", interfacesHttp.AuthMiddleware(interfacesHttp.HandleDeleteBet(useCases)))

	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
	}

	handler := interfacesHttp.CorsMiddlewareDynamic(mux.ServeHTTP, allowedOrigins)

	fmt.Println("---------------------------------------------------")
	fmt.Printf("🚀 Backend (Render-Ready) rodando em http://0.0.0.0:%s\n", port)
	fmt.Println("---------------------------------------------------")

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handler))
}
