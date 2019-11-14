package usecase

import (
	"context"
	"time"

	"github.com/dandynaufaldi/sample-rest-api/author"
	"github.com/dandynaufaldi/sample-rest-api/models"
)

type authorUsecase struct {
	// TODO: add aarticle repository
	authorRepo       author.Repository
	timeoutThreshold time.Duration
}

func (a *authorUsecase) GetByID(ctx context.Context, authorID int64) (*models.Author, error) {
	c, cancel := context.WithTimeout(ctx, a.timeoutThreshold)
	defer cancel()

	author, err := a.authorRepo.GetByID(c, authorID)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (a *authorUsecase) FetchArticle(ctx context.Context, authorID int64, cursor int64, limit int64) ([]*models.Article, error) {

	// c, cancel := context.WithTimeout(ctx, a.timeoutThreshold)
	// defer cancel()

	return nil, nil
}
