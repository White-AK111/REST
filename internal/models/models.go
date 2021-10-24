package models

import (
	"time"
)

// Task structure it's a model for Task entity.
type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

// Repository interface for all repository methods.
type Repository interface {
	CreateTask(text string, tags []string, due time.Time) int
	GetTask(id int) (Task, error)
	DeleteTask(id int) error
	DeleteAllTasks() error
	GetAllTasks() []Task
	GetTasksByTag(tag string) []Task
	GetTasksByDueDate(year int, month time.Month, day int) []Task
}
