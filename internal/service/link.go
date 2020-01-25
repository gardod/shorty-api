package service

import (
	"context"

	"github.com/gardod/shorty-api/internal/driver/redis"
	"github.com/gardod/shorty-api/internal/middleware"
	m "github.com/gardod/shorty-api/internal/model"
	r "github.com/gardod/shorty-api/internal/repository"
	"github.com/sirupsen/logrus"
)

type Link struct {
	log      logrus.FieldLogger
	cache    *redis.Client
	linkRepo *r.Link
}

func NewLink(ctx context.Context) *Link {
	return &Link{
		log:      middleware.GetLogger(ctx),
		cache:    middleware.GetCache(ctx),
		linkRepo: r.NewLink(ctx),
	}
}

func (s *Link) GetByShort(ctx context.Context, short string) (*m.Link, error) {
	link, err := s.linkRepo.Where(r.LinkWhereShort{short}).GetOne(ctx)
	if err != nil {
		return nil, err
	}

	return link, nil
}
