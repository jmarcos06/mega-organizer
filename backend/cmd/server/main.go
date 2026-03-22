package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	appBet "mega-play/internal/application/bet"
	"mega-play/internal/infrastructure/storage"
	interfacesHttp "mega-play/internal/interfaces/http"
)

func main() {
	port := flag.String("port", "8080", "Porta do servidor")
	flag.Parse()

	// Initialize Repository
	repo := storage.NewJSONRepository("apostas_db.json")
	if err := repo.Load(); err != nil {
		log.Printf("⚠️ Aviso: Não foi possível carregar DB existente ou novo: %v", err)
	} else {
		bets, _ := repo.GetAll()
		log.Printf("📂 %d apostas carregadas do histórico.", len(bets))
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

	// Apply CORS middleware to the mux
	handler := interfacesHttp.CorsMiddleware(mux.ServeHTTP)

	// Start Server
	fmt.Println("---------------------------------------------------")
	fmt.Printf("🚀 Backend rodando em http://0.0.0.0:%s\n", *port)
	fmt.Println("---------------------------------------------------")

	log.Fatal(http.ListenAndServe("0.0.0.0:"+*port, handler))
}
