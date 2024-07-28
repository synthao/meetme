package main

import (
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/synthao/meetme/internal/config"
	"log"
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
	stmt, err := db.Prepare("INSERT INTO customers (first_name, last_name, email, phone, address, city, state, zip, country) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	q := `INSERT INTO customers (first_name, last_name, email, phone, address, city, state, zip, country) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	for i := 0; i < *recordsCount; i++ {
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		email := fmt.Sprintf("%s.%s@example.com", firstName, lastName)
		phone := gofakeit.Phone()
		address := gofakeit.Address().Address
		city := gofakeit.City()
		state := gofakeit.State()
		zip := gofakeit.Zip()
		country := gofakeit.Country()

		_, err := db.Exec(q, firstName, lastName, email, phone, address, city, state, zip, country)
		if err != nil {
			log.Fatalf("Failed to insert record: %v", err)
		}

		//customerID, err := exec.LastInsertId()
		//if err != nil {
		//	log.Println(err.Error())
		//	continue
		//}

		//var customerID int64
		//err = db.QueryRow("SELECT lastval()").Scan(&customerID)
		//if err != nil {
		//	log.Fatalf("Failed to get customer ID: %v", err)
		//}

		//orderCount := rand.Intn(3) + 1
		//for j := 0; j < orderCount; j++ {
		//	orderDate := gofakeit.Date()
		//	totalAmount := rand.Float64() * 400
		//	paymentMethod := gofakeit.CreditCardType()
		//	shippingAddress := gofakeit.Address().Address
		//	billingAddress := gofakeit.Address().Address
		//	status := orderStatuses[rand.Intn(len(orderStatuses))]
		//
		//	sqlInsertOrder := `INSERT INTO orders (
		//            				customer_id,
		//            				order_date,
		//            				total_amount,
		//            				status,
		//            				payment_method,
		//            				shipping_address,
		//            				billing_address)
		//						VALUES ($1, $2, $3, $4, $5, $6, $7)`
		//
		//	_, err := db.Exec(sqlInsertOrder, customerID, orderDate, totalAmount, status, paymentMethod, shippingAddress, billingAddress)
		//	if err != nil {
		//		log.Fatalf("Failed to insert order record: %s", err.Error())
		//	}
		//}
	}

	fmt.Println("Records created successfully!")
}
