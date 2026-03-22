package bet

type Bet struct {
	ID         string  `json:"id" bson:"_id"`
	Data       string  `json:"data" bson:"data"`
	Nickname   string  `json:"nickname" bson:"nickname"`
	Season     string  `json:"season" bson:"season"`
	Numeros    []int   `json:"numeros" bson:"numeros"`
	QtdDezenas int     `json:"qtd_dezenas" bson:"qtd_dezenas"`
	Custo      float64 `json:"custo" bson:"custo"`
	Seed       int64   `json:"seed" bson:"seed"`
	Tipo       string  `json:"tipo" bson:"tipo"`
}

type StatNum struct {
	Numero int `json:"numero"`
	Qtd    int `json:"qtd"`
}

type DashboardStats struct {
	TotalGasto     float64   `json:"total_gasto"`
	TotalJogos     int       `json:"total_jogos"`
	UltimasApostas []Bet     `json:"ultimas_apostas"`
	NumerosQuentes []StatNum `json:"numeros_quentes"`
	Seasons        []string  `json:"seasons"`
	Apostadores    []string  `json:"apostadores"`
}

type UserHistory struct {
	Nickname string `json:"nickname"`
	Apostas  []Bet  `json:"apostas"`
}

type Repository interface {
	Save(bet Bet) error
	Delete(id string) error
	GetAll() ([]Bet, error)
	GetByNickname(nickname string) ([]Bet, error)
	CheckCollision(numeros []int) (bool, string, error)
	Load() error
}
