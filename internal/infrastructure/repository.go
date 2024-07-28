package infrastructure

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/synthao/meetme/internal/user/domain"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) domain.Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *domain.User) (int, error) {
	q := "INSERT INTO users (firstname, lastname, email, gender, birth_date) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int
	err := r.db.Get(&id, q, user.FirstName, user.LastName, user.Email, user.Gender, user.BirthDate)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repository) GetByID(id int) (*domain.User, error) {
	var (
		user = &domain.User{}
		dest = &dto{}
		q    = "SELECT id, firstname, lastname, email FROM users WHERE id = $1"
	)

	if err := r.db.Get(dest, q, id); err != nil {
		return nil, err
	}

	user.ID = dest.ID
	user.Email = dest.Email
	user.FirstName = dest.FirstName
	user.LastName = dest.LastName

	return user, nil
}

func (r *repository) GetByEmail(email string) (*domain.User, error) {
	var (
		user = &domain.User{}
		dest = &dto{}
		q    = "SELECT id, firstname, lastname, email FROM users WHERE email = $1"
	)

	if err := r.db.Get(dest, q, email); err != nil {
		return nil, err
	}

	user.ID = dest.ID
	user.Email = dest.Email
	user.FirstName = dest.FirstName
	user.LastName = dest.LastName

	return user, nil
}

func (r *repository) Update(user *domain.User) error {
	q := "UPDATE users SET firstname = $1, lastname = $2, email = $3, updated_at = now() WHERE id = $4"

	_, err := r.db.Exec(q, user.FirstName, user.LastName, user.Email, user.ID)

	return err
}

func (r *repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)

	return err
}

func (r *repository) GetList(limit int, offset int) ([]domain.User, error) {
	//q := "SELECT id, firstname, lastname, email FROM users ORDER BY id LIMIT $1 OFFSET $2"
	q := "SELECT id, firstname, lastname, email FROM users ORDER BY updated_at DESC LIMIT $1 OFFSET $2"

	var dest []dto

	err := r.db.Select(&dest, q, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, len(dest))

	for i, u := range dest {
		users[i] = domain.User{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		}
	}

	return users, nil
}
