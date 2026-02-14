package entities

import "time"

type Link struct {
	Short     string    `db:"short_url"`
	Original  string    `db:"original_url"`
	CreatedAt time.Time `db:"created_at"`
}
