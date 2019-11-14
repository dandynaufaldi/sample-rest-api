package author

import (
	"context"

	"github.com/dandynaufaldi/sample-rest-api/models"
)

// Usecase represent usecase interface
type Usecase interface {
	GetByID(ctx context.Context, authorID int64) (*models.Author, error)
	FetchArticle(ctx context.Context, authorID int64, cursor int64, limit int64) ([]*models.Article, error)
}
