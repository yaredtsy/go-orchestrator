package manager

import (
	"cube/task"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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
		msg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Printf(msg)
		w.WriteHeader(400)
		e := ErrResponse{
			HttpStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
	}

	a.Manager.AddTask(te)
	log.Panicf("Added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}
func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(a.Manager.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskId")
	if taskID == "" {
		log.Printf("No taskId passes in request.\n")
		w.WriteHeader(400)
	}
	tID, _ := uuid.Parse(taskID)
	taskToStop, ok := a.Manager.TaskDb[tID.String()]
	if !ok {
		log.Printf("No task with ID %v found", tID)
		w.WriteHeader(404)
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Completed,
		TimeStamp: time.Now(),
	}

	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	te.Task = taskCopy
	a.Manager.AddTask(te)

	log.Printf("Added task event %v to stop task /%v\n", te.ID)
	w.WriteHeader(204)
}
