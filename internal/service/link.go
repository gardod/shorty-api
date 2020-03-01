package service

import (
	"context"
	"time"

	"github.com/gardod/shorty-api/internal/driver/redis"
	"github.com/gardod/shorty-api/internal/middleware"
	m "github.com/gardod/shorty-api/internal/model"
	r "github.com/gardod/shorty-api/internal/repository"

	vld "github.com/go-ozzo/ozzo-validation/v4"
	vldis "github.com/go-ozzo/ozzo-validation/v4/is"
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

func (s *Link) Validate(link *m.Link) error {
	return vld.ValidateStruct(link,
		vld.Field(&link.Short, vld.Required),
		vld.Field(&link.Long, vld.Required, vldis.URL),
	)
}

func (s *Link) Insert(ctx context.Context, link *m.Link) error {
	// TODO: generate short if not specified

	if err := s.Validate(link); err != nil {
		return err
	}

	return s.linkRepo.Insert(ctx, link)
}

func (s *Link) Update(ctx context.Context, link *m.Link) error {
	if err := s.Validate(link); err != nil {
		return err
	}

	if err := s.cache.Del("link|short:" + link.Short); err != nil {
		s.log.WithError(err).Warning("Unable to delete Link to cache")
	}

	return s.linkRepo.Update(ctx, link)
}

func (s *Link) Delete(ctx context.Context, link *m.Link) error {
	if err := s.cache.Del("link|short:" + link.Short); err != nil {
		s.log.WithError(err).Warning("Unable to delete Link to cache")
	}

	return s.linkRepo.Delete(ctx, link)
}

func (s *Link) Get(ctx context.Context, from time.Time, limit int) ([]m.Link, error) {
	return s.linkRepo.
		Where(r.LinkWhereCreatedAtBefore(from)).
		Order(r.LinkOrderCreatedAt(r.OrderDirectionDesc)).
		WithDeleted(true).
		Get(ctx)
}

func (s *Link) GetByShort(ctx context.Context, short string) (*m.Link, error) {
	var link *m.Link

	key := "link|short:" + short
	if err := s.cache.Get(key, &link); err == nil {
		return link, nil
	} else if err != redis.ErrNotFound {
		s.log.WithError(err).Warning("Unable to get Link from cache")
	}

	link, err := s.linkRepo.Where(r.LinkWhereShort{short}).GetOne(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(key, link, time.Hour); err != nil {
		s.log.WithError(err).Warning("Unable to set Link to cache")
	}

	return link, nil
}

func (s *Link) GetByID(ctx context.Context, id int64) (*m.Link, error) {
	return s.linkRepo.Where(r.LinkWhereID{id}).WithDeleted(true).GetOne(ctx)
}
