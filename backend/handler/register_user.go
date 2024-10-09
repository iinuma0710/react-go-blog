package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iinuma0710/react-go-blog/backend/entity"
)

type RegisterUser struct {
	Service   RegisterUserService
	Validator *validator.Validate
}

func (ru *RegisterUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// リクエストの JSON のでコードに使う構造体
	var b struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}

	// JSON のデコードに失敗した場合の処理
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// JSON のデコード結果を検証してエラーが見つかった場合の処理
	if err := ru.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// データベースにユーザを登録して、エラーがなければ割り振られたユーザ ID を返す
	u, err := ru.Service.RegisterUser(ctx, b.Name, b.Password, b.Role)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID entity.UserID `json:"id"`
	}{ID: u.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
