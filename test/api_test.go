package test

import (
	"ApiBash/handlers"
	"ApiBash/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var db *gorm.DB

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

func TestPost(t *testing.T) {
	setupDb()

	var returnedCommand models.Command

	r := mux.NewRouter()
	r.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
		handlers.CommPostHandler(db, w, r)
	}).Methods("POST")

	command := models.Command{Script: "echo Hello, world!"}
	payload, _ := json.Marshal(command)
	request, _ := http.NewRequest("POST", "/commands", bytes.NewBuffer(payload))
	response := httptest.NewRecorder()

	r.ServeHTTP(response, request)

	if status := response.Code; status != http.StatusCreated {
		t.Errorf("Неверный статус кода: получили %v ожидали %v", status, http.StatusCreated)
	}

	json.NewDecoder(response.Body).Decode(&returnedCommand)
	if returnedCommand.Script != command.Script {
		t.Errorf("Неверный ответ: получили %v ожидали %v", returnedCommand.Script, command.Script)
	}
}

func TestGet(t *testing.T) {
	setupDb()

	r := mux.NewRouter()
	r.HandleFunc("/commands", func(w http.ResponseWriter, r *http.Request) {
		handlers.CommGetHandler(db, w, r)
	}).Methods("GET")

	request, _ := http.NewRequest("GET", "/commands", nil)
	response := httptest.NewRecorder()

	r.ServeHTTP(response, request)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("Неверный статус кода: получили %v ожидали %v", status, http.StatusOK)
	}

	var commands []models.Command
	err := json.Unmarshal(response.Body.Bytes(), &commands)
	if err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}
}

func TestGetId(t *testing.T) {
	setupDb()

	testCommand := models.Command{Script: "Test script"}
	db.Create(&testCommand)

	r := mux.NewRouter()
	r.HandleFunc("/commands/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.CommGetIdHandler(db, w, r)
	}).Methods("GET")

	request, _ := http.NewRequest("GET", fmt.Sprintf("/commands/%d", testCommand.ID), nil)
	response := httptest.NewRecorder()

	r.ServeHTTP(response, request)

	if status := response.Code; status != http.StatusOK {
		t.Errorf("Неверный статус кода: получили %v ожидали %v", status, http.StatusOK)
	}

	var command models.Command
	err := json.Unmarshal(response.Body.Bytes(), &command)
	if err != nil {
		t.Errorf("Ошибка при декодировании ответа: %v", err)
	}

	if command.ID != testCommand.ID {
		t.Errorf("Возвращенная команда не соответствует добавленной")
	}
}
