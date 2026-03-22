package season

import (
	"context"
	"mega-play/internal/domain/season"
)

type UseCases struct {
	repo season.Repository
}

func NewUseCases(repo season.Repository) *UseCases {
	return &UseCases{repo: repo}
}

func (uc *UseCases) CreateSeason(ctx context.Context, name string) error {
	return uc.repo.Save(ctx, season.Season{Name: name})
}

func (uc *UseCases) DeleteSeason(ctx context.Context, name string) error {
	return uc.repo.Delete(ctx, name)
}

func (uc *UseCases) GetAllSeasons(ctx context.Context) ([]season.Season, error) {
	return uc.repo.GetAll(ctx)
}
