package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gardod/shorty-api/internal/driver/postgres"
	"github.com/gardod/shorty-api/internal/middleware"
	"github.com/gardod/shorty-api/internal/model"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Link struct {
	log     logrus.FieldLogger
	db      *postgres.Client
	q       QueryOptions
	withDel bool
}

func NewLink(ctx context.Context) *Link {
	return &Link{
		log: middleware.GetLogger(ctx),
		db:  middleware.GetDB(ctx),
		q:   QueryOptions{},
	}
}

func (r *Link) Insert(ctx context.Context, link *model.Link) error {
	query := `
		INSERT INTO "link" (
			"short",
			"long"
		) VALUES (
			$1,
			$2
		)
		RETURNING
			"created_at",
			"updated_at"
	`

	row := r.db.QueryRowContext(ctx, query, link.Short, link.Long)

	err := row.Scan(&link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		r.log.WithError(err).Error("Unable to insert Link")
		return err
	}

	return nil
}

func (r *Link) Update(ctx context.Context, link *model.Link) error {
	query := `
		UPDATE "link"
		SET "short" = $2,
			"long" = $3,
			"updated_at" = now()
		WHERE "id" = $1
		RETURNING "updated_at"
	`

	row := r.db.QueryRowContext(ctx, query, link.ID, link.Short, link.Long)

	err := row.Scan(&link.UpdatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			r.log.WithError(err).Error("Unable to update Link")
		}
		return err
	}

	return nil
}

func (r *Link) Delete(ctx context.Context, link *model.Link) error {
	query := `
		UPDATE "link"
		SET "deleted_at" = now()
		WHERE "id" = $1
		RETURNING "deleted_at"
	`

	row := r.db.QueryRowContext(ctx, query, link.ID)

	err := row.Scan(&link.DeletedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			r.log.WithError(err).Error("Unable to delete Link")
		}
		return err
	}

	return nil
}

func (r *Link) Get(ctx context.Context) ([]model.Link, error) {
	q := r.q
	if !r.withDel {
		q.Where = append(q.Where, `"l"."deleted_at" IS NULL`)
	}

	query := fmt.Sprintf(`
		SELECT
			%s
			"l"."id",
			"l"."short",
			"l"."long",
			"l"."created_at",
			"l"."updated_at",
			"l"."deleted_at"
		FROM "link" "l"
		%s
		%s
		%s
		%s`,
		q.BuildDistinct(),
		q.BuildJoin(),
		q.BuildWhere(),
		q.BuildOrder(),
		q.BuildLimit(),
	)

	rows, err := r.db.QueryContext(ctx, query, q.Arguments...)
	if err != nil {
		r.log.WithError(err).Error("Unable to select Link")
		return nil, err
	}
	defer rows.Close()

	links := []model.Link{}
	for rows.Next() {
		link := model.Link{}
		err := rows.Scan(
			&link.ID,
			&link.Short,
			&link.Long,
			&link.CreatedAt,
			&link.UpdatedAt,
			&link.DeletedAt,
		)
		if err != nil {
			r.log.WithError(err).Error("Unable to scan Link")
			return nil, err
		}

		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		r.log.WithError(err).Error("Unable to scan Link")
		return nil, err
	}

	return links, nil
}

func (r *Link) GetOne(ctx context.Context) (*model.Link, error) {
	links, err := r.Get(ctx)
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, sql.ErrNoRows
	}

	return &links[0], nil
}

func (r *Link) GetMap(ctx context.Context) (map[int64]model.Link, error) {
	list, err := r.Get(ctx)
	if err != nil {
		return nil, err
	}

	links := make(map[int64]model.Link, len(list))
	for _, link := range list {
		links[link.ID] = link
	}

	return links, nil
}

func (r *Link) GetCount(ctx context.Context) (int, error) {
	q := r.q
	if !r.withDel {
		q.Where = append(q.Where, `"l"."deleted_at" IS NULL`)
	}

	query := fmt.Sprintf(`
		SELECT COUNT("t".*)
		FROM (
			SELECT
				%s
				"l"."id"
			FROM "link" "l"
			%s
			%s
		) "t"`,
		q.BuildDistinct(),
		q.BuildJoin(),
		q.BuildWhere(),
	)

	row := r.db.QueryRowContext(ctx, query, q.Arguments...)

	var count int
	if err := row.Scan(&count); err != nil {
		r.log.WithError(err).Error("Unable to select Link")
		return 0, err
	}

	return count, nil
}

func (r *Link) Distinct(options ...DistinctOption) *Link {
	q := r.q
	q.AddDistinct(options...)
	return &Link{r.log, r.db, q, r.withDel}
}

func (r *Link) Where(options ...WhereOption) *Link {
	q := r.q
	q.AddWhere(options...)
	return &Link{r.log, r.db, q, r.withDel}
}

func (r *Link) Order(options ...OrderOption) *Link {
	q := r.q
	q.AddOrder(options...)
	return &Link{r.log, r.db, q, r.withDel}
}

func (r *Link) Limit(limit int) *Link {
	q := r.q
	q.SetLimit(limit)
	return &Link{r.log, r.db, q, r.withDel}
}

func (r *Link) Offset(offset int) *Link {
	q := r.q
	q.SetOffset(offset)
	return &Link{r.log, r.db, q, r.withDel}
}

func (r *Link) WithDeleted(include bool) *Link {
	return &Link{r.log, r.db, r.q, include}
}

type LinkWhereID []int64

func (o LinkWhereID) GetWhere(start int) (stmt string, args []interface{}) {
	stmt = `"l"."id" = ANY($` + strconv.Itoa(start) + `)`
	args = append(args, pq.Int64Array(o))
	return
}

type LinkWhereShort []string

func (o LinkWhereShort) GetWhere(start int) (stmt string, args []interface{}) {
	stmt = `"l"."short" = ANY($` + strconv.Itoa(start) + `)`
	args = append(args, pq.StringArray(o))
	return
}

type LinkWhereCreatedAtBefore time.Time

func (o LinkWhereCreatedAtBefore) GetWhere(start int) (stmt string, args []interface{}) {
	stmt = `"l"."created_at" < $` + strconv.Itoa(start)
	args = append(args, time.Time(o))
	return
}

type LinkOrderCreatedAt OrderDirection

func (o LinkOrderCreatedAt) GetOrder() string {
	return `"l"."created_at" ` + string(o)
}
