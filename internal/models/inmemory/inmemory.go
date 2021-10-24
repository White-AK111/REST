// Package inmemory provides a simple in-memory "data store" for tasks.
// Tasks are uniquely identified by numeric IDs.
package inmemory

import (
	"fmt"
	"github.com/White-AK111/REST/internal/models"
	"sync"
	"time"
)

// TaskStore is a simple in-memory database of tasks; TaskStore methods are safe to call concurrently.
type TaskStore struct {
	tasks map[int]models.Task
	sync.Mutex
	nextId int
}

// NewStorage function initialize new in-memory repositories.
func NewStorage() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]models.Task)
	ts.nextId = 1
	return ts
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	task := models.Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

// GetTask retrieves a task from the store, by id. If no such id exists, an error is returned.
func (ts *TaskStore) GetTask(id int) (models.Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return models.Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// DeleteTask deletes the task with the given id. If no such id exists, an error is returned.
func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

// DeleteAllTasks deletes all tasks in the store.
func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]models.Task)
	return nil
}

// GetAllTasks returns all the tasks in the store, in arbitrary order.
func (ts *TaskStore) GetAllTasks() []models.Task {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]models.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

// GetTasksByTag returns all the tasks that have the given tag, in arbitrary order.
func (ts *TaskStore) GetTasksByTag(tag string) []models.Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []models.Task

taskLoop:
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue taskLoop
			}
		}
	}
	return tasks
}

// GetTasksByDueDate returns all the tasks that have the given due date, in arbitrary order.
func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []models.Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []models.Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()
		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
