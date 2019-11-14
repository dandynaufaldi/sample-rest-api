package usecase

import (
	"context"
	"time"

	"github.com/dandynaufaldi/sample-rest-api/models"

	"github.com/dandynaufaldi/sample-rest-api/article"
)

// articleUsecase represent article's usecase
type articleUsecase struct {
	articleRepo      article.Repository
	timeoutThreshold time.Duration
}

// func (a *articleUsecase) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, string, error) {
// func (a *articleUsecase) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, string, error) {
// 	if num == 0 {
// 		num = 10
// 	}

// 	return nil, nil, nil
// }

func (a *articleUsecase) GetBySlug(ctx context.Context, slug string) (*models.Article, error) {
	c, cancel := context.WithTimeout(ctx, a.timeoutThreshold)
	defer cancel()

	article, err := a.articleRepo.GetBySlug(c, slug)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (a *articleUsecase) Update(ctx context.Context, article *models.Article) error {
	c, cancel := context.WithTimeout(ctx, a.timeoutThreshold)
	defer cancel()

	err := a.articleRepo.Update(c, article)
	return err
}
