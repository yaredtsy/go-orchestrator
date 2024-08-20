package main

import (
	"cube/manager"
	"cube/node"
	"cube/task"
	"cube/worker"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

// import (
// 	"fmt"
// 	"log"

// 	"github.com/c9s/goprocinfo/linux"
// )

// func main() {
// 	// Read CPU info from /proc/stat
// 	stat, err := linux.ReadStat("/proc/stat")
// 	if err != nil {
// 		log.Fatalf("failed to read /proc/stat: %v", err)
// 	}

// 	for _, s := range stat.CPUStats {
// 		fmt.Printf("CPU: %v, User: %v, System: %v, Idle: %v\n", s.Id, s.User, s.System, s.Idle)
// 	}

// 	// Read Memory info from /proc/meminfo
// 	meminfo, err := linux.ReadMemInfo("/proc/meminfo")
// 	if err != nil {
// 		log.Fatalf("failed to read /proc/meminfo: %v", err)
// 	}

// 	fmt.Printf("MemTotal: %v, MemFree: %v, Buffers: %v, Cached: %v\n",
// 		meminfo.MemTotal, meminfo.MemFree, meminfo.Buffers, meminfo.Cached)

// }

func main() {
	t := task.Task{
		ID:     uuid.New(),
		Name:   "task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}

	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		TimeStamp: time.Now(),
		Task:      t,
	}

	fmt.Println("task: %v\n", t)
	fmt.Printf("task event: %v\n", te)

	w := worker.Worker{
		Name:  "worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	fmt.Printf("worker: %v\n", w)
	w.CollectStats()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]task.Task),
		EventDb: make(map[string][]task.TaskEvent),
		Workers: []string{w.Name},
	}

	fmt.Printf("manger: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	fmt.Printf("node: %v\n", n)

	fmt.Printf("create a test container\n")

	dockerTask, createResult := createContainer()
	if createResult == nil {
		fmt.Printf("%v", createResult.Error)
		os.Exit(1)
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("stop the test container %s\n", createResult.ContainerId)
	_ = stopContainer(dockerTask, createResult.ContainerId)
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
