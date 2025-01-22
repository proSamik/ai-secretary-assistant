package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"
    "github.com/gorilla/mux"
    "github.com/proSamik/aiSecretary/server/models"
    "github.com/proSamik/aiSecretary/server/repository"
)

var (
    todoRepo *repository.TodoRepository
    once     sync.Once
)

func initRepository() {
    once.Do(func() {
        todoRepo = repository.NewTodoRepository()
    })
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
    initRepository()

    var input models.TodoInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todo, err := todoRepo.Create(&input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
    initRepository()

    status := r.URL.Query().Get("status")
    
    todos, err := todoRepo.List(status)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
    initRepository()

    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    todo, err := todoRepo.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
    initRepository()

    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var input struct {
        Status string `json:"status"`
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todo, err := todoRepo.Update(id, input.Status)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
    initRepository()

    params := mux.Vars(r)
    id, err := strconv.Atoi(params["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := todoRepo.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}