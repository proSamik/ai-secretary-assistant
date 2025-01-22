package models

import (
    "time"
)

type Todo struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    DueDate     time.Time `json:"due_date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type TodoInput struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    DueDate     string `json:"due_date"`
}