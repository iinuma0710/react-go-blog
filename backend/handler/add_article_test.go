package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/store"
	"github.com/iinuma0710/react-go-blog/backend/testutil"
)

func TestAddArticle(t *testing.T) {
	t.Parallel()

	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_article/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_article/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_article/bad_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_article/bad_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/articles",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			sut := AddArticle{
				Store: &store.ArticleStore{
					Articles: map[entity.ArticleID]*entity.Article{},
				},
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			rsp := w.Result()
			testutil.AssertResponse(t, rsp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
