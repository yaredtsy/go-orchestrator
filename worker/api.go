package worker

import (
	"cube/task"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
}

type ErrResponse struct {
	HttpStatusCode int
	Message        string
}

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}

	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error decoding task event: %v", err)
		log.Print(msg)
		w.WriteHeader(http.StatusBadRequest)

		e := ErrResponse{
			HttpStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	a.Worker.AddTask(te.Task)
	log.Printf("Added task %v\n", te.Task.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(te)
}

func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	if taskID == "" {
		log.Print("No task ID provided")
		w.WriteHeader(400)
	}

	tID, _ := uuid.Parse(taskID)
	_, ok := a.Worker.Db[tID]
	if !ok {
		log.Printf("No task found with ID %v", tID)
		w.WriteHeader(404)
	}

	taskToStop := a.Worker.Db[tID]
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)

	log.Printf("Task %v stopped in container%v\n", tID, taskCopy.ContainerID)
	w.WriteHeader(204)
}

func (a *Api) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Worker.Stats)
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})
	a.Router.Route("/task", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
}

func (a *Api) Start() {
	a.initRouter()
	log.Printf("Starting API server on %s:%d\n", a.Address, a.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
}
