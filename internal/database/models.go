// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID
	CreatedAt     sql.NullTime
	UpdatedAt     sql.NullTime
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
}

type FeedFollow struct {
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type User struct {
	ID        uuid.UUID
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Name      string
}
