package main

import (
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/synthao/meetme/internal/config"
	"github.com/synthao/meetme/internal/user/domain"
	"log"
	"math/rand"
	"time"
)

var orderStatuses = []string{
	"pending",
	"processing",
	"shipped",
	"delivered",
	"cancelled",
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	recordsCount := flag.Int("count", 100, "Number of customer records to create")

	flag.Parse()

	// Подключение к базе данных PostgreSQL
	db, err := sqlx.Connect("postgres", config.GetDSN())
	if err != nil {
		log.Fatalf("Could not connect to the database: %v, dsn: %s", err, config.GetDSN())
	}
	defer db.Close()

	// Создание подготовленного запроса
	q := "INSERT INTO users (firstname, lastname, email, gender, created_at, birth_date) VALUES ($1, $2, $3, $4, now(), $5)"

	for i := 0; i < *recordsCount; i++ {
		person := gofakeit.Person()
		firstName := person.FirstName
		lastName := person.LastName
		email := gofakeit.Email()

		_, err := db.Exec(q, firstName, lastName, email, gender(person.Gender), generateBirthDate())
		if err != nil {
			log.Fatalf("Failed to insert record: %v", err)
		}
	}

	fmt.Println("Records created successfully!")
}

func gender(g string) domain.Gender {
	if g == "male" {
		return domain.GenderMale
	}
	return domain.GenderFemale
}

func generateBirthDate() time.Time {
	now := time.Now()
	minAge := 18
	maxAge := 100 // Можно настроить максимальный возраст по вашему усмотрению

	// Вычисляем минимальную и максимальную даты рождения
	minBirthDate := now.AddDate(-(minAge + 1), 0, 0)
	maxBirthDate := now.AddDate(-maxAge, 0, 0)

	// Вычисляем разницу в днях между минимальной и максимальной датами
	diff := minBirthDate.Sub(maxBirthDate).Hours() / 24

	// Генерируем случайное количество дней в пределах разницы
	randomDays := rand.Float64() * diff

	// Вычисляем случайную дату рождения
	birthDate := maxBirthDate.AddDate(0, 0, int(randomDays))

	return birthDate
}
