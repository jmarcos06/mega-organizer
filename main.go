package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --- Configurações e Constantes ---

const (
	PrecoApostaSimples = 6.00
	ArquivoDB          = "apostas_db.json" // Mudamos para JSON para facilitar estruturas complexas
	MaxNumeros         = 20
	MinNumeros         = 6
	TotalDezenasJogo   = 60
	AccessToken        = "MONOBOLA123"
)

// --- Estruturas de Dados ---

type Aposta struct {
	ID         string  `json:"id"`
	Data       string  `json:"data"`
	Nickname   string  `json:"nickname"` // Novo: Quem apostou
	Season     string  `json:"season"`   // Nova: Temporada/Categoria
	Numeros    []int   `json:"numeros"`
	QtdDezenas int     `json:"qtd_dezenas"`
	Custo      float64 `json:"custo"`
	Seed       int64   `json:"seed"`
	Tipo       string  `json:"tipo"`
}

// DataStore gerencia o acesso seguro ao arquivo
type DataStore struct {
	mu      sync.RWMutex
	Apostas []Aposta `json:"apostas"`
}

// Estrutura para resposta da API (AJAX)
type APIResponse struct {
	TotalGasto     float64   `json:"total_gasto"`
	TotalJogos     int       `json:"total_jogos"`
	UltimasApostas []Aposta  `json:"ultimas_apostas"`
	NumerosQuentes []StatNum `json:"numeros_quentes"`
	Seasons        []string  `json:"seasons"`
}

// Resposta para histórico do usuário
type UserHistoryResponse struct {
	Nickname string   `json:"nickname"`
	Apostas  []Aposta `json:"apostas"`
}

type StatNum struct {
	Numero int `json:"numero"`
	Qtd    int `json:"qtd"`
}

var store = DataStore{
	Apostas: []Aposta{},
}

// --- Lógica de Negócio (Core) ---

func CalcularCusto(qtdDezenas int) float64 {
	if qtdDezenas < 6 {
		return 0
	}
	combs := new(big.Int).Binomial(int64(qtdDezenas), 6)
	return float64(combs.Int64()) * PrecoApostaSimples
}

func GerarJogo(qtd int, fixos []int, seed int64) ([]int, error) {
	if qtd < MinNumeros || qtd > MaxNumeros {
		return nil, fmt.Errorf("qtd inválida")
	}

	var r *rand.Rand
	if seed != 0 {
		r = rand.New(rand.NewSource(seed))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	escolhidos := make(map[int]bool)
	resultado := make([]int, 0, qtd)

	for _, num := range fixos {
		if !escolhidos[num] {
			escolhidos[num] = true
			resultado = append(resultado, num)
		}
	}

	for len(resultado) < qtd {
		num := r.Intn(TotalDezenasJogo) + 1
		if !escolhidos[num] {
			escolhidos[num] = true
			resultado = append(resultado, num)
		}
	}

	sort.Ints(resultado)
	return resultado, nil
}

// --- Persistência (JSON com Mutex) ---

func (ds *DataStore) Carregar() error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	file, err := os.Open(ArquivoDB)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &ds.Apostas)
}

func (ds *DataStore) Salvar() error {
	ds.mu.RLock() // Lock de leitura para serializar
	data, err := json.MarshalIndent(ds.Apostas, "", "  ")
	ds.mu.RUnlock()

	if err != nil {
		return err
	}

	ds.mu.Lock() // Lock de escrita para salvar no disco
	defer ds.mu.Unlock()
	return os.WriteFile(ArquivoDB, data, 0644)
}

func (ds *DataStore) Adicionar(a Aposta) error {
	ds.mu.Lock()
	ds.Apostas = append(ds.Apostas, a)
	ds.mu.Unlock()
	return ds.Salvar()
}

// Verifica se existe colisão exata de jogo
func (ds *DataStore) VerificarColisao(numeros []int) (bool, string) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	// Cria string representativa para comparação simples
	novoStr := fmt.Sprint(numeros)

	for _, a := range ds.Apostas {
		existenteStr := fmt.Sprint(a.Numeros)
		if novoStr == existenteStr {
			return true, a.Nickname
		}
	}
	return false, ""
}

// Estatísticas para o Dashboard
func (ds *DataStore) ObterStats(currentSeason string) APIResponse {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	totalGasto := 0.0
	freqMap := make(map[int]int)
	seasonsMap := make(map[string]bool)

	// Copia para evitar race conditions na ordenação
	copiaApostas := make([]Aposta, 0)
	for _, a := range ds.Apostas {
		// Migrate old apostas without season
		if a.Season == "" {
			a.Season = "Mega da Virada 2025"
		}
		seasonsMap[a.Season] = true

		// Filter by current season if specified
		if currentSeason == "" || a.Season == currentSeason {
			copiaApostas = append(copiaApostas, a)
		}
	}

	// Ordena por data decrescente
	sort.Slice(copiaApostas, func(i, j int) bool {
		return copiaApostas[i].Data > copiaApostas[j].Data
	})

	for _, a := range copiaApostas {
		totalGasto += a.Custo
		for _, n := range a.Numeros {
			freqMap[n]++
		}
	}

	// Top Números
	var stats []StatNum
	for k, v := range freqMap {
		stats = append(stats, StatNum{Numero: k, Qtd: v})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Qtd > stats[j].Qtd
	})
	if len(stats) > 5 {
		stats = stats[:5]
	}

	// Convert seasons map to slice
	var seasons []string
	for season := range seasonsMap {
		seasons = append(seasons, season)
	}
	sort.Strings(seasons)

	return APIResponse{
		TotalGasto:     totalGasto,
		TotalJogos:     len(copiaApostas),
		UltimasApostas: copiaApostas,
		NumerosQuentes: stats,
		Seasons:        seasons,
	}
}

// Obter histórico de um usuário
func (ds *DataStore) ObterHistoricoUsuario(nickname string) UserHistoryResponse {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var apostasUsuario []Aposta
	for _, a := range ds.Apostas {
		if strings.EqualFold(a.Nickname, nickname) {
			apostasUsuario = append(apostasUsuario, a)
		}
	}

	// Ordena por data decrescente
	sort.Slice(apostasUsuario, func(i, j int) bool {
		return apostasUsuario[i].Data > apostasUsuario[j].Data
	})

	return UserHistoryResponse{
		Nickname: nickname,
		Apostas:  apostasUsuario,
	}
}

// Remover aposta por ID
func (ds *DataStore) RemoverAposta(id string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	for i, a := range ds.Apostas {
		if a.ID == id {
			// Remove elemento do slice
			ds.Apostas = append(ds.Apostas[:i], ds.Apostas[i+1:]...)
			return ds.Salvar()
		}
	}
	return fmt.Errorf("aposta não encontrada")
}

// --- Interface Web & Handlers ---

// --- Handlers HTTP ---

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleGetData(w http.ResponseWriter, r *http.Request) {
	currentSeason := r.URL.Query().Get("season")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.ObterStats(currentSeason))
}

func handleGetCost(w http.ResponseWriter, r *http.Request) {
	qtdStr := r.URL.Query().Get("qtd")
	qtd, err := strconv.Atoi(qtdStr)
	if err != nil || qtd < 6 || qtd > 20 {
		http.Error(w, `{"error": "Quantidade inválida"}`, http.StatusBadRequest)
		return
	}

	custo := CalcularCusto(qtd)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"custo": custo})
}

func handleGetUserHistory(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	if nickname == "" {
		http.Error(w, `{"error": "Nickname é obrigatório"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(store.ObterHistoricoUsuario(nickname))
}

func handleDeleteBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "ID é obrigatório"}`, http.StatusBadRequest)
		return
	}

	if err := store.RemoverAposta(id); err != nil {
		http.Error(w, `{"error": "Aposta não encontrada"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handlePostBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Estrutura do payload recebido
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

	// Conversão
	qtd, _ := strconv.Atoi(input.QtdStr)
	var fixos []int
	if input.FixosStr != "" {
		// Support both comma and space separation
		input.FixosStr = strings.ReplaceAll(input.FixosStr, ",", " ")
		parts := strings.Fields(input.FixosStr)
		for _, p := range parts {
			val, err := strconv.Atoi(strings.TrimSpace(p))
			if err == nil && val >= 1 && val <= 60 {
				fixos = append(fixos, val)
			}
		}
	}

	// Validação básica
	if input.Nickname == "" {
		input.Nickname = "Anônimo"
	}
	if input.Season == "" {
		input.Season = "Mega da Virada 2025"
	}

	// Gera o jogo
	seed := time.Now().UnixNano()
	nums, err := GerarJogo(qtd, fixos, seed)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Verifica Colisão
	colisao, quem := store.VerificarColisao(nums)

	// Cria Objeto
	custo := CalcularCusto(qtd)
	tipo := "Simples"
	if qtd > 6 {
		tipo = "Desdobramento"
	}

	aposta := Aposta{
		ID:         fmt.Sprintf("%x", seed%1000000),
		Data:       time.Now().Format(time.RFC3339),
		Nickname:   input.Nickname,
		Season:     input.Season,
		Numeros:    nums,
		QtdDezenas: qtd,
		Custo:      custo,
		Seed:       seed,
		Tipo:       tipo,
	}

	if err := store.Adicionar(aposta); err != nil {
		http.Error(w, `{"error": "Erro ao salvar"}`, http.StatusInternalServerError)
		return
	}

	// Resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"aposta":      aposta,
		"colisao":     colisao,
		"colisao_com": quem,
	})
}

// --- Middleware de Autenticação Simples ---
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token") // ou: r.Header.Get("Authorization")
		if token != AccessToken {
			http.Error(w, "Acesso não autorizado", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// --- Main ---

func main() {
	// Flags CLI mantidas para uso local rápido, se desejado
	cliMode := flag.Bool("cli", false, "Modo CLI (sem servidor)")
	port := flag.String("port", "8080", "Porta do servidor")
	flag.Parse()

	// Carrega banco de dados
	if err := store.Carregar(); err != nil {
		log.Printf("⚠️  Aviso: Não foi possível carregar DB existente ou novo: %v", err)
	} else {
		log.Printf("📂 %d apostas carregadas do histórico.", len(store.Apostas))
	}

	if *cliMode {
		fmt.Println("Modo CLI desativado nesta versão servidor. Use a interface web.")
		return
	}

	// Rotas
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/dados", handleGetData)
	http.HandleFunc("/api/apostar", handlePostBet)
	http.HandleFunc("/api/custo", handleGetCost)
	http.HandleFunc("/api/usuario/historico", handleGetUserHistory)
	http.HandleFunc("/api/aposta/deletar", AuthMiddleware(handleDeleteBet))

	// Inicia Server
	fmt.Println("---------------------------------------------------")
	fmt.Printf("🚀 Mega Hub rodando em http://0.0.0.0:%s\n", *port)
	fmt.Println("💻 Acesse do navegador. Compartilhe seu IP com amigos.")
	fmt.Println("---------------------------------------------------")

	log.Printf("Servidor ouvindo na porta %s...", *port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+*port, nil))
}
