package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dandynaufaldi/sample-rest-api/author/mock"
	"github.com/golang/mock/gomock"

	"github.com/dandynaufaldi/sample-rest-api/models"
)

func Test_authorUsecase_GetByID(t *testing.T) {
	type args struct {
		ctx      context.Context
		authorID int64
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthorRepository := mock.NewMockRepository(ctrl)
	mockAuthor := &models.Author{
		ID:   1,
		Name: "Dandy Naufaldi",
	}
	mockAuthorRepository.EXPECT().GetByID(gomock.Any(), 1).Return(mockAuthor, nil).Times(1)
	// field := fields{
	// 	authorRepo:       mockAuthorRepository,
	// 	timeoutThreshold: 10,
	// }
	tests := []struct {
		name    string
		args    args
		want    *models.Author
		wantErr bool
	}{
		// TODO: Add test cases.
		// name: "Success",
		// args: args{
		// 	ctx: context.Background(),
		// 	0
		// }

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authorUsecase{
				authorRepo:       mockAuthorRepository,
				timeoutThreshold: 100 * time.Millisecond,
			}
			got, err := a.GetByID(tt.args.ctx, tt.args.authorID)

			if (err != nil) != tt.wantErr {
				t.Errorf("authorUsecase.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authorUsecase.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
