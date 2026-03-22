package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"mega-play/internal/domain/bet"
)

type JSONRepository struct {
	mu      sync.RWMutex
	arquivo string
	apostas []bet.Bet
}

func NewJSONRepository(arquivo string) *JSONRepository {
	return &JSONRepository{
		arquivo: arquivo,
		apostas: []bet.Bet{},
	}
}

func (r *JSONRepository) Load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(r.arquivo)
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

	return json.Unmarshal(bytes, &r.apostas)
}

func (r *JSONRepository) saveInternal() error {
	data, err := json.MarshalIndent(r.apostas, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.arquivo, data, 0644)
}

func (r *JSONRepository) Save(b bet.Bet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.apostas = append(r.apostas, b)
	return r.saveInternal()
}

func (r *JSONRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, a := range r.apostas {
		if a.ID == id {
			r.apostas = append(r.apostas[:i], r.apostas[i+1:]...)
			return r.saveInternal()
		}
	}
	return fmt.Errorf("aposta não encontrada")
}

func (r *JSONRepository) GetAll() ([]bet.Bet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copia := make([]bet.Bet, len(r.apostas))
	copy(copia, r.apostas)
	return copia, nil
}

func (r *JSONRepository) GetByNickname(nickname string) ([]bet.Bet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []bet.Bet
	for _, a := range r.apostas {
		if strings.EqualFold(a.Nickname, nickname) {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *JSONRepository) CheckCollision(numeros []int) (bool, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	novoStr := fmt.Sprint(numeros)

	for _, a := range r.apostas {
		existenteStr := fmt.Sprint(a.Numeros)
		if novoStr == existenteStr {
			return true, a.Nickname, nil
		}
	}
	return false, "", nil
}
