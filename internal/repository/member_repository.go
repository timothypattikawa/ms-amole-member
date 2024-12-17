package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/timothypattikawa/ms-kamoro-costumer/internal/repository/postgres"
	"time"
)

type MemberRepository interface {
	ExecTx(ctx context.Context, fn func(q *sqlc.Queries) error) error
	Exec(ctx context.Context, fn func(q *sqlc.Queries) error) error
}

type MemberRepositoryImpl struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func (m MemberRepositoryImpl) ExecTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	begin, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	tx := sqlc.New(m.db).WithTx(begin)

	err = fn(tx)
	if err != nil {
		err := begin.Rollback(ctx)
		if err != nil {
			return err
		}
	}

	err = begin.Commit(ctx)
	if err != nil {
		return err
	}

	return err
}

func (m MemberRepositoryImpl) Exec(ctx context.Context, fn func(q *sqlc.Queries) error) error {

	err := fn(m.queries)
	if err != nil {
		return err
	}

	return err
}

func NewMemberRepository(db *pgxpool.Pool) MemberRepository {
	return &MemberRepositoryImpl{
		db:      db,
		queries: sqlc.New(db),
	}
}
