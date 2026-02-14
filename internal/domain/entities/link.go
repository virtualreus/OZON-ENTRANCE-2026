package entities

import "time"

type Link struct {
	Short     string    `db:"short"`
	Original  string    `db:"original"`
	CreatedAt time.Time `db:"created_at"`
}
