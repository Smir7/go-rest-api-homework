package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTask(w http.ResponseWriter, r *http.Request) {
	var taskArray []Task
	for _, task := range tasks {
		taskArray = append(taskArray, task)
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getIdTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]

	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if ok {
		http.Error(w, "Уже есть такая задача", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	resp, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	w.Header().Del(task.ID)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	r.Get("/tasks", getTask)
	r.Get("/tasks/{id}", getIdTask)
	r.Post("/tasks", postTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
