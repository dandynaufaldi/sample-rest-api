package author

import (
	"context"

	"github.com/dandynaufaldi/sample-rest-api/models"
)

// Repository represent repository interface
type Repository interface {
	GetByID(ctx context.Context, authorID int64) (*models.Author, error)
}
