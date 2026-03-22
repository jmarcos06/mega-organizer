package season

import "context"

type Season struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

type Repository interface {
	Save(ctx context.Context, s Season) error
	Delete(ctx context.Context, name string) error
	GetAll(ctx context.Context) ([]Season, error)
}
