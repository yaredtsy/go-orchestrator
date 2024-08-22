package main

import (
	"cube/task"
	"cube/worker"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	fmt.Println("Starting Cube worker")
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}
	go runTasks(&w)
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently. \n")

		}
		log.Println("Sleeping for 5 seconds")
		time.Sleep(10 * time.Second)
	}
}

func createContainer() (*task.Docker, *task.DockerResult) {
	c := task.Config{
		Name:  "test-container-1",
		Image: "postgres:16",
		Env:   []string{"POSTGRES_PASSWORD=secret", "POSTGRES_USER=cube"},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv)
	d := task.Docker{
		Client: dc,
		Config: c,
	}

	result := d.Run()
	if result.Error != nil {
		fmt.Printf("Error running container: %v\n", result.Error)
		return nil, nil
	}

	fmt.Printf("Container %s started\n", result.ContainerId, c)
	return &d, &result

}

func stopContainer(d *task.Docker, id string) *task.DockerResult {
	result := d.Stop(id)

	if result.Error != nil {
		fmt.Printf("Error stopping container: %v\n", result.Error)
		return nil
	}

	fmt.Printf("Container %s stopped and removed\n", id)
	return &result
}
