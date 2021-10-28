package gin_gonic

import (
	"github.com/White-AK111/REST/config"
	"github.com/White-AK111/REST/internal/models"
	"github.com/White-AK111/REST/internal/models/inmemory"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// getAllTasksHandler handler for GET method without id.
func (ts *taskServer) getAllTasksHandler(c *gin.Context) {
	allTasks := ts.store.GetAllTasks()
	c.JSON(http.StatusOK, allTasks)
}

// deleteAllTasksHandler handler for DELETE method without id.
func (ts *taskServer) deleteAllTasksHandler(c *gin.Context) {
	err := ts.store.DeleteAllTasks()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
	}
}

// createTaskHandler handler for POST method do create task.
func (ts *taskServer) createTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}

	var rt RequestTask
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

// getTaskHandler handler for GET method with id.
func (ts *taskServer) getTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

// deleteTaskHandler handler for DELETE method with id.
func (ts *taskServer) deleteTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if err = ts.store.DeleteTask(id); err != nil {
		c.String(http.StatusNotFound, err.Error())
	}
}

// tagHandler handler for "tag" path.
func (ts *taskServer) tagHandler(c *gin.Context) {
	tag := c.Params.ByName("tag")
	tasks := ts.store.GetTasksByTag(tag)
	c.JSON(http.StatusOK, tasks)
}

// dueHandler handler for "due" path.
func (ts *taskServer) dueHandler(c *gin.Context) {
	badRequestError := func() {
		c.String(http.StatusBadRequest, "expect /due/<year>/<month>/<day>, got %v", c.FullPath())
	}

	year, err := strconv.Atoi(c.Params.ByName("year"))
	if err != nil {
		badRequestError()
		return
	}

	month, err := strconv.Atoi(c.Params.ByName("month"))
	if err != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	day, err := strconv.Atoi(c.Params.ByName("day"))
	if err != nil {
		badRequestError()
		return
	}

	tasks := ts.store.GetTasksByDueDate(year, time.Month(month), day)
	c.JSON(http.StatusOK, tasks)
}

// Init function do initialize a new server with parameters from config.yaml.
func Init(cfg *config.Config) {
	router := gin.Default()
	var server *taskServer

	switch cfg.Server.TypeOfRepository {
	case "in-memory":
		server = NewTaskServerInmemory()
	default:
		cfg.ErrorLogger.Fatal("Unknown repository type.")
	}

	router.POST("/task/", server.createTaskHandler)
	router.GET("/task/", server.getAllTasksHandler)
	router.DELETE("/task/", server.deleteAllTasksHandler)
	router.GET("/task/:id", server.getTaskHandler)
	router.DELETE("/task/:id", server.deleteTaskHandler)
	router.GET("/tag/:tag", server.tagHandler)
	router.GET("/due/:year/:month/:day", server.dueHandler)

	cfg.ErrorLogger.Printf("Start server %s with storage %s on: %s\n", cfg.Server.TypeOfServer, cfg.Server.TypeOfRepository, cfg.Server.ServerAddress+":"+strconv.Itoa(cfg.Server.ServerPort))
	err := router.Run(cfg.Server.ServerAddress + ":" + strconv.Itoa(cfg.Server.ServerPort))
	if err != nil {
		cfg.ErrorLogger.Fatalf("Error on start server: %s\n", err)
	}
}
