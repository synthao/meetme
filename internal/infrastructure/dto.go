package infrastructure

type dto struct {
	ID        int    `db:"id"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	Email     string `db:"email"`
}
