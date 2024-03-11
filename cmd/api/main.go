package main

import (
	"ApiBash/handlers"
	"ApiBash/models"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB

func startServer(db *gorm.DB) {
	r := mux.NewRouter()

	r.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
		handlers.CommPostHandler(db, w, r)
	}).Methods("POST")
	r.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
		handlers.CommGetHandler(db, w, r)
	}).Methods("GET")
	r.HandleFunc("/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.CommGetIdHandler(db, w, r)
	}).Methods("GET")

	err := http.ListenAndServe(":6000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func setupDb() *gorm.DB {
	var err error

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Command{})

	return db
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден!")
	}

	db = setupDb()

	fmt.Println("Запуск сервера...")
	startServer(db)
}
