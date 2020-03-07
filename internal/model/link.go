package model

import (
	"regexp"
	"time"

	vld "github.com/go-ozzo/ozzo-validation/v4"
	vldis "github.com/go-ozzo/ozzo-validation/v4/is"
)

type LinkRequest struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

type Link struct {
	LinkRequest
	ID        int64      `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (l *Link) Validate() error {
	return vld.ValidateStruct(l,
		vld.Field(&l.Short,
			vld.Required,
			vld.Match(regexp.MustCompile("^[a-zA-Z0-9-]+$")).Error("can only contain letters, digits, and hyphens"),
		),
		vld.Field(&l.Long,
			vld.Required,
			vldis.URL,
		),
	)
}
