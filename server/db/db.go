package db

import (
    "database/sql"
    "fmt"
    "os"
    "log"
)

var DB *sql.DB

func InitDB() error {
    if DB != nil {
        return nil
    }

    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    log.Printf("Connecting to database with connection string: %s", connStr)

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error opening database: %v", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("error connecting to the database: %v", err)
    }

    log.Printf("Successfully connected to database")

    // Create todos table if not exists
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS todos (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            status VARCHAR(50) NOT NULL DEFAULT 'pending',
            due_date TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return fmt.Errorf("error creating todos table: %v", err)
    }

    return nil
}