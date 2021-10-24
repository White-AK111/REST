// Package gorilla it's a basic example of a REST server with several routes, using the gorilla/mux
package gorilla

import (
	"encoding/json"
	"fmt"
	"github.com/White-AK111/REST/config"
	"github.com/White-AK111/REST/internal/models"
	"github.com/White-AK111/REST/internal/models/inmemory"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

// createTaskHandler handler for POST method do create task.
func (ts *taskServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

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
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	allTasks := ts.store.GetAllTasks()
	renderJSON(w, allTasks)
}

// getTaskHandler handler for GET method with id.
func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	// Here and elsewhere, not checking error of convert because the router only matches the [0-9]+ regex.
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task, err := ts.store.GetTask(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

// deleteTaskHandler handler for DELETE method with id.
func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete task at %s\n", req.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := ts.store.DeleteTask(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// deleteAllTasksHandler handler for DELETE method without id.
func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all tasks at %s\n", req.URL.Path)
	err := ts.store.DeleteAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// tagHandler handler for "tag" path.
func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by tag at %s\n", req.URL.Path)

	tag := mux.Vars(req)["tag"]
	tasks := ts.store.GetTasksByTag(tag)
	renderJSON(w, tasks)
}

// dueHandler handler for "due" path.
func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling tasks by due at %s\n", req.URL.Path)

	vars := mux.Vars(req)
	badRequestError := func() {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	year, _ := strconv.Atoi(vars["year"])
	month, _ := strconv.Atoi(vars["month"])
	if month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, _ := strconv.Atoi(vars["day"])

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	renderJSON(w, tasks)
}

// Init function do initialize a new server with parameters from config.yaml.
func Init(cfg *config.Config) {
	router := mux.NewRouter()
	router.StrictSlash(true)
	var server *taskServer

	switch cfg.Service.TypeOfRepository {
	case "in-memory":
		server = NewTaskServerInmemory()
	default:
		log.Fatal("Unknown repository type.")
	}

	router.HandleFunc("/task/", server.createTaskHandler).Methods("POST")
	router.HandleFunc("/task/", server.getAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/", server.deleteAllTasksHandler).Methods("DELETE")
	router.HandleFunc("/task/{id:[0-9]+}/", server.getTaskHandler).Methods("GET")
	router.HandleFunc("/task/{id:[0-9]+}/", server.deleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/tag/{tag}/", server.tagHandler).Methods("GET")
	router.HandleFunc("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}/", server.dueHandler).Methods("GET")

	log.Printf("Start service on: %s\n", cfg.Service.ServerAddress+":"+strconv.Itoa(cfg.Service.ServerPort))
	log.Fatal(http.ListenAndServe(cfg.Service.ServerAddress+":"+strconv.Itoa(cfg.Service.ServerPort), router))
}
