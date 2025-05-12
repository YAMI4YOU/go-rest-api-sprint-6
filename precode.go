package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

// Ниже напишите обработчики для каждого эндпоинта

// Обработчик для получения списка всех задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	// сериализируем данные из мапы в JSON для ответа клиенту
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// записываем в заголовок тип ответа
	w.Header().Set("Content-Type", "application/json")
	// записываем статус ответа
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные JSON данные в ответ
	w.Write(resp)
}

// Обработчик добавления задачи при получении JSON
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task

	// десериализируем полученные данные из JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// добавляем задачу в мапу
	tasks[task.ID] = task

	// записываем в заголовок тип ответа
	w.Header().Set("Content-Type", "application/json")
	// записываем статус ответа
	w.WriteHeader(http.StatusCreated)
}

// Обработчик возврата задачи по ID
func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметр url из запроса
	id := chi.URLParam(r, "id")

	// проверяем существует ли ключ в мапе
	task, ok := tasks[id]
	if !ok {
		http.Error(w, `"Задача отсутствует"`, http.StatusBadRequest)
		return
	}

	// сериализируем данные из мапы в JSON для ответа клиенту
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// записываем в заголовок тип ответа
	w.Header().Set("Content-Type", "application/json")
	// записываем статус ответа
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные JSON данные в ответ
	w.Write(resp)
}

// Обработчик удаления задачи по ID
func delTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметр url из запроса
	id := chi.URLParam(r, "id")

	// проверяем есть ли такая задача в мапе
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача отсутствует", http.StatusBadRequest)
		return
	}

	// удаляем задачу из мапы по ID
	delete(tasks, id)

	// записываем в заголовок тип ответа
	w.Header().Set("Content-Type", "application/json")
	// записываем статус ответа
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Логирования отображения методов и статусов в терминале
	r.Use(middleware.Logger)

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasksHandler)
	r.Post("/tasks", createTaskHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Delete("/tasks/{id}", delTaskByIDHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
