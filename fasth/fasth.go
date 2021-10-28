package fasth

import (
	"encoding/json"
	"fmt"
	"github.com/White-AK111/REST/config"
	"github.com/White-AK111/REST/internal/models"
	"github.com/White-AK111/REST/internal/models/inmemory"
	"github.com/White-AK111/REST/middleware"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"mime"
	"net/http"
	"strconv"
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
func renderJSONFast(c *fasthttp.RequestCtx, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	c.SetContentType("application/json")
	_, err = c.Write(js)
	if err != nil {
		c.Error(err.Error(), http.StatusInternalServerError)
		return
	}
}

// getAllTasksHandler handler for GET method without id.
func (ts *taskServer) getAllTasksHandler(c *fasthttp.RequestCtx) {
	allTasks := ts.store.GetAllTasks()
	renderJSONFast(c, allTasks)
}

// deleteAllTasksHandler handler for DELETE method without id.
func (ts *taskServer) deleteAllTasksHandler(c *fasthttp.RequestCtx) {
	err := ts.store.DeleteAllTasks()
	if err != nil {
		c.Error(err.Error(), http.StatusNotFound)
	}
}

// createTaskHandler handler for POST method do create task.
func (ts *taskServer) createTaskHandler(c *fasthttp.RequestCtx) {
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
	contentType := c.Request.Header.Peek("Content-Type")
	mediaType, _, err := mime.ParseMediaType(string(contentType))
	if err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		c.Error("expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	var rt RequestTask
	if err := json.Unmarshal(c.PostBody(), &rt); err != nil {
		c.Error(err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSONFast(c, ResponseId{Id: id})
}

// getTaskHandler handler for GET method with id.
func (ts *taskServer) getTaskHandler(c *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(c.UserValue("id").(string))
	task, err := ts.store.GetTask(id)
	if err != nil {
		c.Error(err.Error(), http.StatusNotFound)
		return
	}

	renderJSONFast(c, task)
}

// deleteTaskHandler handler for DELETE method with id.
func (ts *taskServer) deleteTaskHandler(c *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(c.UserValue("id").(string))
	err := ts.store.DeleteTask(id)
	if err != nil {
		c.Error(err.Error(), http.StatusNotFound)
	}
}

// tagHandler handler for "tag" path.
func (ts *taskServer) tagHandler(c *fasthttp.RequestCtx) {
	tag := c.UserValue("tag").(string)
	tasks := ts.store.GetTasksByTag(tag)
	renderJSONFast(c, tasks)
}

// dueHandler handler for "due" path.
func (ts *taskServer) dueHandler(c *fasthttp.RequestCtx) {
	badRequestError := func() {
		c.Error(fmt.Sprintf("expect /due/<year>/<month>/<day>, got %v", string(c.Path())), http.StatusBadRequest)
	}

	year, _ := strconv.Atoi(c.UserValue("year").(string))
	month, _ := strconv.Atoi(c.UserValue("month").(string))
	if month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}
	day, _ := strconv.Atoi(c.UserValue("day").(string))

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	renderJSONFast(c, tasks)
}

func (ts *taskServer) panicHandler(c *fasthttp.RequestCtx) {
	panic("test panic")
}

// Init function do initialize a new server with parameters from config.yaml.
func Init(cfg *config.Config) {
	r := router.New()
	var server *taskServer

	switch cfg.Server.TypeOfRepository {
	case "in-memory":
		server = NewTaskServerInmemory()
	default:
		cfg.ErrorLogger.Fatal("Unknown repository type.")
	}

	r.POST("/task/", server.createTaskHandler)
	r.GET("/task/", server.getAllTasksHandler)
	r.DELETE("/task/", server.deleteAllTasksHandler)
	r.GET("/task/{id:[0-9]+}", server.getTaskHandler)
	r.DELETE("/task/{id:[0-9]+}", server.deleteTaskHandler)
	r.GET("/tag/{tag}", server.tagHandler)
	r.GET("/due/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", server.dueHandler)
	// For test panic
	r.GET("/panic", server.panicHandler)

	cfg.ErrorLogger.Printf("Start server %s with storage %s on: %s\n", cfg.Server.TypeOfServer, cfg.Server.TypeOfRepository, cfg.Server.ServerAddress+":"+strconv.Itoa(cfg.Server.ServerPort))

	s := &fasthttp.Server{
		Handler: middleware.LoggerAndPanicRecover(r.Handler),
		Name:    "fastHttpWithLoggerAndPanicRecover",
	}

	err := s.ListenAndServe(cfg.Server.ServerAddress + ":" + strconv.Itoa(cfg.Server.ServerPort))
	if err != nil {
		cfg.ErrorLogger.Fatalf("Error on start server: %s\n", err)
	}
}
