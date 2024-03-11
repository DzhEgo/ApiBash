package handlers

import (
	"ApiBash/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os/exec"
)

func CommPostHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var commands models.Command

	err := json.NewDecoder(r.Body).Decode(&commands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer r.Body.Close()

	result := db.Create(&commands)
	if result.Error != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	go func(db *gorm.DB, cmd *models.Command) {
		out, err := exec.Command("sh", "-c", cmd.Script).CombinedOutput()
		cmd.Executed = true
		if err != nil {
			log.Printf("Ошибка выполнения команды: %v. Вывод: %s", err, string(out))
			cmd.Result = err.Error()
		} else {
			cmd.Result = string(out)
		}
		db.Save(cmd)
		fmt.Println(cmd.Result)
		fmt.Println(cmd.Executed)
	}(db, &commands)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commands)
}

func CommGetHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var commands []models.Command
	result := db.Find(&commands)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

func CommGetIdHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var commands models.Command
	result := db.First(&commands, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}
