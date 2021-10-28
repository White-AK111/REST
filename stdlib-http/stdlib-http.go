// Package stdlib_http it's a basic example of a REST server with several routes, using only the standard library.
package stdlib_http

import (
	"encoding/json"
	"fmt"
	"github.com/White-AK111/REST/config"
	"github.com/White-AK111/REST/internal/models"
	"github.com/White-AK111/REST/internal/models/inmemory"
	"github.com/White-AK111/REST/middleware"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// taskServer struct for server of task/
type taskServer struct {
	store models.Repository
}

// NewTaskServerInmemory function initialize a new taskServer.
func NewTaskServerInmemory() *taskServer {
	store := inmemory.NewStorage()
	return &taskServer{store: store}
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// taskHandler handler for "task" path.
func (ts *taskServer) taskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		// Request is plain "/task/", without trailing ID.
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTasksHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ts.deleteAllTasksHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// Request has an ID, as in "/task/<id>".
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.deleteTaskHandler(w, req, id)
		} else if req.Method == http.MethodGet {
			ts.getTaskHandler(w, req, id)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

// createTaskHandler handler for POST method do create task.
func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	// Types used internally in this handler to (de-)serialize the request and response from/to JSON.
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSON(w, ResponseId{Id: id})
}

// getAllTasksHandler handler for GET method without id.
func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := ts.store.GetAllTasks()
	renderJSON(w, allTasks)
}

// getTaskHandler handler for GET method with id.
func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

// deleteTaskHandler handler for DELETE method with id.
func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request, id int) {
	err := ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// deleteAllTasksHandler handler for DELETE method without id.
func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	err := ts.store.DeleteAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// tagHandler handler for "tag" path.
func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /tag/<tag>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /tag/<tag> path", http.StatusBadRequest)
		return
	}
	tag := pathParts[1]

	tasks := ts.store.GetTasksByTag(tag)
	renderJSON(w, tasks)
}

// dueHandler handler for "due" path.
func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /due/<date>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}
	if len(pathParts) != 4 {
		badRequestError()
		return
	}

	year, err := strconv.Atoi(pathParts[1])
	if err != nil {
		badRequestError()
		return
	}
	month, err := strconv.Atoi(pathParts[2])
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, err := strconv.Atoi(pathParts[3])
	if err != nil {
		badRequestError()
		return
	}

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	renderJSON(w, tasks)
}

// Init function do initialize a new server with parameters from config.yaml.
func Init(cfg *config.Config) {
	mux := http.NewServeMux()
	var server *taskServer

	switch cfg.Server.TypeOfRepository {
	case "in-memory":
		server = NewTaskServerInmemory()
	default:
		cfg.ErrorLogger.Fatal("Unknown repository type.")
	}

	mux.HandleFunc("/task/", server.taskHandler)
	mux.HandleFunc("/tag/", server.tagHandler)
	mux.HandleFunc("/due/", server.dueHandler)

	handler := middleware.Logging(mux)
	handler = middleware.PanicRecovery(handler)

	cfg.ErrorLogger.Printf("Start server %s with storage %s on: %s\n", cfg.Server.TypeOfServer, cfg.Server.TypeOfRepository, cfg.Server.ServerAddress+":"+strconv.Itoa(cfg.Server.ServerPort))
	err := http.ListenAndServe(cfg.Server.ServerAddress+":"+strconv.Itoa(cfg.Server.ServerPort), handler)
	if err != nil {
		cfg.ErrorLogger.Fatalf("Error on start server: %s\n", err)
	}
}
