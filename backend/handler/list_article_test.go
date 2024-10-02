package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/testutil"
)

func TestListTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		articles []*entity.Article
		want     want
	}{
		"ok": {
			articles: []*entity.Article{
				{
					ID:     1,
					Title:  "test1",
					Status: entity.ArticlePublished,
				},
				{
					ID:     2,
					Title:  "test2",
					Status: entity.ArticleDraft,
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_article/ok_rsp.json.golden",
			},
		},
		"empty": {
			articles: []*entity.Article{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_article/empty_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			moq := &ListArticlesServiceMock{}
			moq.ListArticlesFunc = func(ctx context.Context) (entity.Articles, error) {
				if tt.articles != nil {
					return tt.articles, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := ListArticle{Service: moq}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
