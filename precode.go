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

// Обработчик для получения всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tasksSlice := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		tasksSlice = append(tasksSlice, task)
	}
	if err := json.NewEncoder(w).Encode(tasksSlice); err != nil {
		http.Error(w, "Ошибка при получении задач", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil || newTask.ID == "" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}
	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при получении задачи", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обработчик удаления задачи по ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
