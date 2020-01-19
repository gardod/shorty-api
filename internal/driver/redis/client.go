package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	redis "github.com/go-redis/redis/v7"
)

type Client struct {
	pool *redis.Client
}

func NewClient(ctx context.Context) *Client {
	if pool == nil {
		return &Client{}
	}
	return &Client{pool: pool.WithContext(ctx)}
}

func (c *Client) Get(key string, result interface{}) error {
	if c.pool == nil {
		return ErrNotFound
	}

	cached, err := c.pool.Get(key).Bytes()
	if err == redis.Nil {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	buf := bytes.NewBuffer(cached)
	enc := gob.NewDecoder(buf)
	err = enc.Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	if c.pool == nil {
		return nil
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	err = c.pool.Set(key, buf.Bytes(), expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Del(keys ...string) error {
	if c.pool == nil {
		return nil
	}
	return c.pool.Del(keys...).Err()
}

func (c *Client) FlushDB() error {
	if c.pool == nil {
		return nil
	}
	return c.pool.FlushDB().Err()
}
