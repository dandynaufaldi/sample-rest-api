package article

import (
	"context"

	"github.com/dandynaufaldi/sample-rest-api/models"
)

// Repository represent article's repository
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, string, error)
	GetBySlug(ctx context.Context, slug string) (*models.Article, error)
	Update(ctx context.Context, article *models.Article) error
	SearchByTitle(ctx context.Context, title string)
	Store(ctx context.Context, article *models.Article) error
	Delete(ctx context.Context, slug string) error
}
