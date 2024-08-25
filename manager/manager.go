package manager

import (
	"bytes"
	"cube/task"
	"cube/worker"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending       queue.Queue
	TaskDb        map[string]*task.Task
	EventDb       map[string]*task.TaskEvent
	Workers       []string
	LastWorker    int
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func New(workers []string) *Manager {
	return &Manager{
		Workers: workers,
	}
}

func (m *Manager) GetTasks() []*task.Task {
	tasks := []*task.Task{}
	for _, t := range m.TaskDb {
		tasks = append(tasks, t)
	}
	return tasks
}

func (m *Manager) AddTask(task task.TaskEvent) {
	m.EventDb[task.ID.String()] = &task
	m.TaskDb[task.Task.ID.String()] = &task.Task
}
func (m *Manager) SelectWorker() string {
	var newWorker int
	if m.LastWorker+1 < len(m.Workers) {
		newWorker = m.LastWorker + 1
		m.LastWorker++
	} else {
		newWorker = 0
		m.LastWorker = 0
	}
	return m.Workers[newWorker]
}

func (m *Manager) ProcessTasks() {
	for {
		log.Printf("Processing any tasks in the queue")
		m.SendWork()
		log.Println("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
	}
}
func (m *Manager) UpdateTasks() {
	for {
		log.Println("Checking for task updates from workers")
		m.updateTasks()
		log.Println("Task updates compeleted")
		log.Println("Sleeping for 15 seconds")
		time.Sleep(15 * time.Second)
	}
}
func (m *Manager) updateTasks() {
	for _, worker := range m.Workers {
		log.Printf("Checking worker %v for rask updates", worker)
		url := fmt.Sprintf("http://%s/task", worker)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error connecting to %v: %v\n", worker, err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("error sending requestL %v\n", err)
		}

		d := json.NewDecoder(resp.Body)
		var tasks []*task.Task
		err = d.Decode(&tasks)
		if err != nil {
			log.Printf("Error unmarshalling tasks: %s\n", err.Error())
		}
	}
}

func (m *Manager) SendWork() {
	if m.Pending.Len() > 0 {
		w := m.SelectWorker()

		e := m.Pending.Dequeue()
		te := e.(task.TaskEvent)
		t := te.Task
		log.Printf("Pulled %v off pending queue\n", t)

		m.EventDb[t.ID.String()] = &te
		m.WorkerTaskMap[w] = append(m.WorkerTaskMap[w], t.ID)
		m.TaskWorkerMap[t.ID] = w

		t.State = task.Scheduled
		m.TaskDb[t.ID.String()] = &t

		data, err := json.Marshal(te)
		if err != nil {
			log.Printf("Unable to marshal object: %v.\n", t)
		}

		url := fmt.Sprintf("http://%s/task", w)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error connecting to %v: %v", w, err)
			m.Pending.Enqueue(t)
			return
		}

		d := json.NewDecoder(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			e := worker.ErrResponse{}
			err := d.Decode(&e)
			if err != nil {
				fmt.Printf("Error decoding response: %s\n", err.Error())
				return
			}
			log.Printf("Response error (%d): %s", e.HttpStatusCode, e.Message)
			return
		}

		t = task.Task{}
		err = d.Decode(&t)
		if err != nil {
			fmt.Printf("Error decoding response: %s\n", err.Error())
			return
		}
		log.Printf("%#v\n", t)
	} else {
		log.Println("No work in the queue")
	}

}
