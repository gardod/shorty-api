package model

import "time"

type Link struct {
	ID        int64      `json:"id"`
	Short     string     `json:"short"`
	Long      string     `json:"long"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
