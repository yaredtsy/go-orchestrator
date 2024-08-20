package main

import (
	"cube/task"
	"cube/worker"
	"fmt"
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
	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-task-2",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("starting task")
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerId

	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 30)

	fmt.Printf("stopping task %s\n", t.ID)
	t.State = task.Completed
	w.AddTask(t)
	result = w.RunTask()
	if result.Error != nil {
		panic(result.Error)
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
