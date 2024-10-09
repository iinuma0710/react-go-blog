package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/iinuma0710/react-go-blog/backend/entity"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	// レコードの作成時刻と更新時刻を取得
	u.CreatedAt = r.Clocker.Now()
	u.UpdatedAt = r.Clocker.Now()

	// SQL の発行と実行
	sql := `INSERT INTO user (
		name, password, role, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}

	// 登録されたレコードの ID を取得
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}
