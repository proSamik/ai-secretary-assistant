package repository

import (
    "database/sql"
    "time"
    "fmt"
    "github.com/proSamik/aiSecretary/server/db"
    "github.com/proSamik/aiSecretary/server/models"
)

type TodoRepository struct {
    db *sql.DB
}

var todoRepo *TodoRepository

func NewTodoRepository() *TodoRepository {
    if todoRepo != nil {
        return todoRepo
    }
    
    if db.DB == nil {
        panic("Database connection not initialized")
    }
    
    todoRepo = &TodoRepository{db: db.DB}
    return todoRepo
}

func (r *TodoRepository) Create(todo *models.TodoInput) (*models.Todo, error) {
    var newTodo models.Todo
    
    if r.db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }
    
    dueDate, err := time.Parse("2006-01-02", todo.DueDate)
    if err != nil {
        return nil, fmt.Errorf("invalid date format: %v", err)
    }

    err = r.db.QueryRow(`
        INSERT INTO todos (title, description, status, due_date)
        VALUES ($1, $2, $3, $4)
        RETURNING id, title, description, status, due_date, created_at, updated_at
    `, todo.Title, todo.Description, "pending", dueDate).Scan(
        &newTodo.ID, &newTodo.Title, &newTodo.Description, 
        &newTodo.Status, &newTodo.DueDate, &newTodo.CreatedAt, &newTodo.UpdatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("error creating todo: %v", err)
    }

    return &newTodo, nil
}

func (r *TodoRepository) List(status string) ([]models.Todo, error) {
    if r.db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }

    var query string
    var rows *sql.Rows
    var err error

    if status != "" {
        query = `SELECT * FROM todos WHERE status = $1 ORDER BY created_at DESC`
        rows, err = r.db.Query(query, status)
    } else {
        query = `SELECT * FROM todos ORDER BY created_at DESC`
        rows, err = r.db.Query(query)
    }

    if err != nil {
        return nil, fmt.Errorf("error querying todos: %v", err)
    }
    defer rows.Close()

    var todos []models.Todo
    for rows.Next() {
        var todo models.Todo
        err := rows.Scan(
            &todo.ID, &todo.Title, &todo.Description, &todo.Status,
            &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning todo row: %v", err)
        }
        todos = append(todos, todo)
    }

    return todos, nil
}

func (r *TodoRepository) GetByID(id int) (*models.Todo, error) {
    if r.db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }

    var todo models.Todo
    err := r.db.QueryRow("SELECT * FROM todos WHERE id = $1", id).Scan(
        &todo.ID, &todo.Title, &todo.Description, &todo.Status,
        &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("todo not found")
    }
    if err != nil {
        return nil, fmt.Errorf("error getting todo: %v", err)
    }
    return &todo, nil
}

func (r *TodoRepository) Update(id int, status string) (*models.Todo, error) {
    if r.db == nil {
        return nil, fmt.Errorf("database connection is nil")
    }

    var todo models.Todo
    err := r.db.QueryRow(`
        UPDATE todos 
        SET status = $1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $2
        RETURNING *
    `, status, id).Scan(
        &todo.ID, &todo.Title, &todo.Description, &todo.Status,
        &todo.DueDate, &todo.CreatedAt, &todo.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("todo not found")
    }
    if err != nil {
        return nil, fmt.Errorf("error updating todo: %v", err)
    }
    return &todo, nil
}

func (r *TodoRepository) Delete(id int) error {
    if r.db == nil {
        return fmt.Errorf("database connection is nil")
    }

    result, err := r.db.Exec("DELETE FROM todos WHERE id = $1", id)
    if err != nil {
        return fmt.Errorf("error deleting todo: %v", err)
    }

    count, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error getting affected rows: %v", err)
    }

    if count == 0 {
        return fmt.Errorf("todo not found")
    }

    return nil
}