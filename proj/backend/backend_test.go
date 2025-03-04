package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCreateTodo(t *testing.T) {
	getDB()

	newTodo := map[string]string{"title": "Test Task"}
	body, _ := json.Marshal(newTodo)

	req, err := http.NewRequest("POST", "/createTodo", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleCreate)

	handler.ServeHTTP(rr, req)
	fmt.Println(rr)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 201 Created, got %v", rr.Code)
	}
}

func TestCreateTodoDB(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	pool = &pgxpool.Pool{}

	mock.ExpectExec("INSERT INTO todos").
		WithArgs("Test Task").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = createTodo("Test Task")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
