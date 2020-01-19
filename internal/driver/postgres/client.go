package postgres

import (
	"context"
	"database/sql"
)

type Client struct {
	pool *sql.DB
	tx   *sql.Tx
}

func NewClient() *Client {
	return &Client{pool: pool}
}

func (c *Client) Begin(ctx context.Context) error {
	if c.tx != nil {
		return ErrTxActive
	}

	tx, err := c.pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	c.tx = tx
	return nil
}

func (c *Client) Commit() error {
	if c.tx == nil {
		return ErrTxInactive
	}

	err := c.tx.Commit()
	if err != nil {
		return err
	}

	c.tx = nil
	return nil
}

func (c *Client) Rollback() error {
	if c.tx == nil {
		return ErrTxInactive
	}

	err := c.tx.Rollback()
	if err != nil {
		return err
	}

	c.tx = nil
	return nil
}

func (c *Client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if c.tx != nil {
		return c.tx.ExecContext(ctx, query, args...)
	}
	return c.pool.ExecContext(ctx, query, args...)
}

func (c *Client) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if c.tx != nil {
		return c.tx.PrepareContext(ctx, query)
	}
	return c.pool.PrepareContext(ctx, query)
}

func (c *Client) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if c.tx != nil {
		return c.tx.QueryContext(ctx, query, args...)
	}
	return c.pool.QueryContext(ctx, query, args...)
}

func (c *Client) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if c.tx != nil {
		return c.tx.QueryRowContext(ctx, query, args...)
	}
	return c.pool.QueryRowContext(ctx, query, args...)
}
