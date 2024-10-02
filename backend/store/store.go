package store

import (
	"errors"

	"github.com/iinuma0710/react-go-blog/backend/entity"
)

type ArticleStore struct {
	// 動作確認のための仮実装なので、あえてエクスポート
	LastID   entity.ArticleID
	Articles map[entity.ArticleID]*entity.Article
}

var (
	Articles    = &ArticleStore{Articles: map[entity.ArticleID]*entity.Article{}}
	ErrNotFound = errors.New("not found")
)

// 記事を一つ追加する
func (as *ArticleStore) Add(a *entity.Article) (entity.ArticleID, error) {
	as.LastID++
	a.ID = as.LastID
	as.Articles[a.ID] = a
	return a.ID, nil
}

// ソート済みのタスク一覧を返す
func (as *ArticleStore) All() entity.Articles {
	articles := make([]*entity.Article, len(as.Articles))
	for i, a := range as.Articles {
		articles[i-1] = a
	}
	return articles
}
