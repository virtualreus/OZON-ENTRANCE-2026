package entities

type Link struct {
	Short     string `db:"short"`
	Original  string `db:"original"`
	CreatedAt int64  `db:"created_at"`
}
