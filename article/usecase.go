package article

import (
	"context"

	"github.com/dandynaufaldi/sample-rest-api/models"
)

// Usecase represent article's usecase
type Usecase interface {
	Fetch(ctx context.Context, cursor string, limit int64) ([]*models.Article, string, error)
	// FetchByTitle(ctx context.Context, title string, cursor int64, limit int64) ([]*models.Article, string, error)
	// FetchByAuthor(ctx context.Context, authorID int64, cursor int64, limit int64) ([]*models.Article, string, error)
	GetByTitle(ctx context.Context, title string) (*models.Article, error)
	GetBySlug(ctx context.Context, slug string) (*models.Article, error)
	Update(ctx context.Context, article *models.Article) error
	Store(ctx context.Context, article *models.Article) error
	Delete(ctx context.Context, slug string) error
}
