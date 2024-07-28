package domain

type Gender int

const (
	GenderMale Gender = iota
	GenderFemale
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	BirthDate string
	Gender    Gender
}

type Repository interface {
	Create(user *User) (int, error)
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Delete(id int) error
	GetList(limit int, offset int) ([]User, error)
	Update(user *User) error
}
