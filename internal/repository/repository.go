package repository

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("resource already exists")
)

// Repository handles data persistence.
// Replace the in-memory store with your database of choice (PostgreSQL, MySQL, etc.)
type Repository struct {
	mu    sync.RWMutex
	users map[string]*User
}

func New() *Repository {
	return &Repository{
		users: make(map[string]*User),
	}
}

// For database connections, you would typically:
//
//	func New(db *sql.DB) *Repository {
//	    return &Repository{db: db}
//	}
//
// Or with GORM:
//
//	func New(db *gorm.DB) *Repository {
//	    return &Repository{db: db}
//	}
