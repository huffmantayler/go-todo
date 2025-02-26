package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func (t *Todo) printTodo() {
	fmt.Printf("ID: %d, Title: %s, Done: %v\n", t.ID, t.Title, t.Done)
}

var (
	pool *pgxpool.Pool
	once sync.Once
)

type TitleRequest struct {
	Title string `json:"title"`
}

func main() {
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	mux.HandleFunc("/createTodo", handleCreate)
	mux.HandleFunc("/deleteTodo", handleDelete)
	mux.HandleFunc("/updateTodo", handleUpdate)
	mux.HandleFunc("/getAllTodos", handleGetAll)

	defer closeDB()

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}

func handleGetAll(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	todos, err := getAllTodos(db)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error getting all todos")
		fmt.Println(err)
		http.Error(w, "error getting all todos", http.StatusBadRequest)
		return
	}

	for _, todo := range todos {
		todo.printTodo()
	}

	jsonErr := json.NewEncoder(w).Encode(todos)
	if jsonErr != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	title := r.URL.Query().Get("title")
	doneStr := r.URL.Query().Get("done")

	fmt.Printf("ID: %s, Title: %s, Done: %s\n", idStr, title, doneStr)

	// Check if at least one of the parameters is present
	if (title != "" && doneStr != "") || (title == "" && doneStr == "") {
		http.Error(w, "Please provide either 'title' or 'done', but not both", http.StatusBadRequest)
		return
	}

	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter, must be an integer", http.StatusBadRequest)
		return
	}
	db := getDB()

	updateTodo(id, doneStr, title, db)
	io.WriteString(w, "Todo updated\n")
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	var titleRequest TitleRequest
	err := json.NewDecoder(r.Body).Decode(&titleRequest)
	// todo.printTodo()
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db := getDB()

	createTodo(titleRequest.Title, db)
	io.WriteString(w, "Todo created\n")
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	idToDeleteStr := r.URL.Query().Get("id")

	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method, expected DELETE", http.StatusMethodNotAllowed)
		return
	}

	if idToDeleteStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	idToDelete, err := strconv.Atoi(idToDeleteStr)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter, must be an integer", http.StatusBadRequest)
		return
	}

	db := getDB()
	deleteTodo(idToDelete, db)

	io.WriteString(w, "Todo Deleted\n")
}

func getDB() *pgxpool.Pool {
	once.Do(func() {
		cfg := LoadConfig()
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
		var err error
		pool, err = pgxpool.New(context.Background(), connStr)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		fmt.Println("Connected to PostgreSQL")
	})
	return pool
}

func closeDB() {
	if pool != nil {
		pool.Close()
		fmt.Println("Database connection pool closed.")
	}
}

func createTodo(title string, pool *pgxpool.Pool) error {

	query := `INSERT INTO todos (title) VALUES ($1) RETURNING id`

	var id int

	err := pool.QueryRow(context.Background(), query, title).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	fmt.Printf("Created todo with ID: %d\n", id)
	return nil
}

func updateTodo(id int, doneStr string, title string, pool *pgxpool.Pool) error {
	var query string
	var err error

	if doneStr != "" {
		var done bool
		done, err = strconv.ParseBool(doneStr)
		if err != nil {
			return fmt.Errorf("invalid 'done' parameter, must be a boolean")
		}
		query = `UPDATE todos SET done = $1 WHERE id = $2`
		_, err = pool.Exec(context.Background(), query, done, id)
	} else if title != "" {
		query = `UPDATE todos SET title = $1 WHERE id = $2`
		_, err = pool.Exec(context.Background(), query, title, id)
	} else {
		return fmt.Errorf("neither 'done' nor 'title' parameter provided")
	}

	if err != nil {
		return fmt.Errorf("error updating todo: %w", err)
	}

	fmt.Printf("Todo with ID %d updated\n", id)
	return nil
}

func deleteTodo(id int, pool *pgxpool.Pool) error {

	query := `DELETE FROM todos WHERE id = $1`

	_, err := pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error deleting todo: %w", err)
	}

	fmt.Printf("Todo with ID %d deleted\n", id)
	return nil
}

func getAllTodos(pool *pgxpool.Pool) ([]Todo, error) {

	var todos []Todo

	query := `SELECT id, title, done FROM todos`

	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error getting all todos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var done bool
		err := rows.Scan(&id, &title, &done)
		if err != nil {
			return nil, fmt.Errorf("error scanning todo: %w", err)
		}
		todos = append(todos, Todo{id, title, done})
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", rows.Err())
	}

	return todos, nil
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .evn file")
	}
	return Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

}
