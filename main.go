package main

import (
	"cube/manager"
	"cube/task"
	"cube/worker"
	"fmt"
	"os"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	whost := os.Getenv("CUBE_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	fmt.Println("Starting Cube worker")
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	wapi := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go wapi.Start()

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)
	mApi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.UpdateTasks()
	go m.ProcessTasks()
	mApi.Start()

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
