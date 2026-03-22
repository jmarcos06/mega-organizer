package bet

import (
	"fmt"
	"sort"
	"time"

	domainBet "mega-play/internal/domain/bet"
)

type UseCases struct {
	repo domainBet.Repository
}

func NewUseCases(r domainBet.Repository) *UseCases {
	return &UseCases{repo: r}
}

func (u *UseCases) CreateBet(nickname string, qtd int, fixos []int, season string) (domainBet.Bet, bool, string, error) {
	if nickname == "" {
		nickname = "Anônimo"
	}
	if season == "" {
		season = "Mega da Virada 2025"
	}

	seed := time.Now().UnixNano()
	nums, err := domainBet.GerarJogo(qtd, fixos, seed)
	if err != nil {
		return domainBet.Bet{}, false, "", err
	}

	colisao, quem, _ := u.repo.CheckCollision(nums)

	custo := domainBet.CalcularCusto(qtd)
	tipo := "Simples"
	if qtd > 6 {
		tipo = "Desdobramento"
	}

	bet := domainBet.Bet{
		ID:         fmt.Sprintf("%x", seed%1000000),
		Data:       time.Now().Format(time.RFC3339),
		Nickname:   nickname,
		Season:     season,
		Numeros:    nums,
		QtdDezenas: qtd,
		Custo:      custo,
		Seed:       seed,
		Tipo:       tipo,
	}

	if err := u.repo.Save(bet); err != nil {
		return domainBet.Bet{}, false, "", err
	}

	return bet, colisao, quem, nil
}

func (u *UseCases) DeleteBet(id string) error {
	return u.repo.Delete(id)
}

func (u *UseCases) GetUserHistory(nickname string) (domainBet.UserHistory, error) {
	apostas, err := u.repo.GetByNickname(nickname)
	if err != nil {
		return domainBet.UserHistory{}, err
	}

	sort.Slice(apostas, func(i, j int) bool {
		return apostas[i].Data > apostas[j].Data
	})

	return domainBet.UserHistory{
		Nickname: nickname,
		Apostas:  apostas,
	}, nil
}

func (u *UseCases) GetDashboardStats(currentSeason string) (domainBet.DashboardStats, error) {
	todas, err := u.repo.GetAll()
	if err != nil {
		return domainBet.DashboardStats{}, err
	}

	totalGasto := 0.0
	freqMap := make(map[int]int)
	seasonsMap := make(map[string]bool)
	apostadoresMap := make(map[string]bool)

	filtradas := make([]domainBet.Bet, 0)
	for _, a := range todas {
		if a.Season == "" {
			a.Season = "Mega da Virada 2025"
		}
		seasonsMap[a.Season] = true

		if currentSeason == "" || a.Season == currentSeason {
			filtradas = append(filtradas, a)
		}
	}

	sort.Slice(filtradas, func(i, j int) bool {
		return filtradas[i].Data > filtradas[j].Data
	})

	for _, a := range filtradas {
		totalGasto += a.Custo
		apostadoresMap[a.Nickname] = true
		for _, n := range a.Numeros {
			freqMap[n]++
		}
	}

	var stats []domainBet.StatNum
	for k, v := range freqMap {
		stats = append(stats, domainBet.StatNum{Numero: k, Qtd: v})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Qtd > stats[j].Qtd
	})
	if len(stats) > 5 {
		stats = stats[:5]
	}

	var seasons []string
	for season := range seasonsMap {
		seasons = append(seasons, season)
	}
	sort.Strings(seasons)

	var apostadores []string
	for p := range apostadoresMap {
		apostadores = append(apostadores, p)
	}
	sort.Strings(apostadores)

	return domainBet.DashboardStats{
		TotalGasto:     totalGasto,
		TotalJogos:     len(filtradas),
		UltimasApostas: filtradas,
		NumerosQuentes: stats,
		Seasons:        seasons,
		Apostadores:    apostadores,
	}, nil
}
