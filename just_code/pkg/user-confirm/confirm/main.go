package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/confirmUser", func(w http.ResponseWriter, r *http.Request) {
		confirmUserHandler(w, r, db) // Передача db в обработчик
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func confirmUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Error(w, "Только POST запросы разрешены", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow("SELECT id FROM user_confirmations WHERE code = $1", request.Code).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		http.Error(w, "Код не найден", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("Ошибка при запросе в базу данных: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("UPDATE user_confirmations SET confirmed = TRUE WHERE id = $1", id)
	if err != nil {
		log.Printf("Ошибка при обновлении пользователя: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	json.NewEncoder(w).Encode(struct{ Status string }{"Пользователь подтвержден"})
}
