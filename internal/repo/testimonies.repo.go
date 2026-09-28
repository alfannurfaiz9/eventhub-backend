package repo

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TestimonyRepo struct {
	db *pgxpool.Pool
}

func NewTestimonyRepo(db *pgxpool.Pool) *TestimonyRepo {
	return &TestimonyRepo{
		db: db,
	}
}

func (t *TestimonyRepo) GetTestimony(ctx context.Context) ([]model.TestimonyList, error) {
	sql := `
	SELECT  testimonies.id, users.full_name, testimonies.company, testimonies.position, testimonies.message
	FROM testimonies
	JOIN users ON users.id = testimonies.user_id`

	rows, err := t.db.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	var testimonyLists []model.TestimonyList
	for rows.Next() {
		var testimonyList model.TestimonyList

		if err := rows.Scan(
			&testimonyList.Testimony.Id,
			&testimonyList.User.FullName,
			&testimonyList.Testimony.Company,
			&testimonyList.Testimony.Position,
			&testimonyList.Testimony.Message,
		); err != nil {
			return nil, err
		}

		testimonyLists = append(testimonyLists, testimonyList)
	}

	return testimonyLists, nil
}

func (t *TestimonyRepo) SetTestimony(ctx context.Context, userId int, body model.Testimony) error {
	sql := `
	INSERT INTO testimonies(user_id, message, company, position)
	VALUES($1, $2, $3, $4)`
	args := []any{userId, body.Message, body.Company, body.Position}

	cmd, err := t.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
