package bet

import (
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"time"
)

const (
	PrecoApostaSimples = 6.00
	MaxNumeros         = 20
	MinNumeros         = 6
	TotalDezenasJogo   = 60
)

// CalcularCusto calculates the cost of a bet based on number of tens
func CalcularCusto(qtdDezenas int) float64 {
	if qtdDezenas < MinNumeros {
		return 0
	}
	combs := new(big.Int).Binomial(int64(qtdDezenas), 6)
	return float64(combs.Int64()) * PrecoApostaSimples
}

// GerarJogo generates a random game allowing some fixed numbers
func GerarJogo(qtd int, fixos []int, seed int64) ([]int, error) {
	if qtd < MinNumeros || qtd > MaxNumeros {
		return nil, fmt.Errorf("quantidade inválida")
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
