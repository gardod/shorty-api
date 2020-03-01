package model

import "time"

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
