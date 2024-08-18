package main

import (
	"fmt"
	"log"

	"github.com/c9s/goprocinfo/linux"
)

func main() {
	// Read CPU info from /proc/stat
	stat, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Fatalf("failed to read /proc/stat: %v", err)
	}

	for _, s := range stat.CPUStats {
		fmt.Printf("CPU: %v, User: %v, System: %v, Idle: %v\n", s.Id, s.User, s.System, s.Idle)
	}

	// Read Memory info from /proc/meminfo
	meminfo, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatalf("failed to read /proc/meminfo: %v", err)
	}

	fmt.Printf("MemTotal: %v, MemFree: %v, Buffers: %v, Cached: %v\n",
		meminfo.MemTotal, meminfo.MemFree, meminfo.Buffers, meminfo.Cached)

}

// func main() {
// 	t := task.Task{
// 		ID:     uuid.New(),
// 		Name:   "task-1",
// 		State:  task.Pending,
// 		Image:  "Image-1",
// 		Memory: 1024,
// 		Disk:   1,
// 	}

// 	te := task.TaskEvent{
// 		ID:        uuid.New(),
// 		State:     task.Pending,
// 		TimeStamp: time.Now(),
// 		Task:      t,
// 	}

// 	fmt.Println("task: %v\n", t)
// 	fmt.Printf("task event: %v\n", te)

// 	w := worker.Worker{
// 		Name:  "worker-1",
// 		Queue: *queue.New(),
// 		Db:    make(map[uuid.UUID]*task.Task),
// 	}
// 	fmt.Printf("worker: %v\n", w)
// 	w.CollectStats()
// 	w.StartTask()
// 	w.StopTask()

// 	m := manager.Manager{
// 		Pending: *queue.New(),
// 		TaskDb:  make(map[string][]task.Task),
// 		EventDb: make(map[string][]task.TaskEvent),
// 		Workers: []string{w.Name},
// 	}

// 	fmt.Printf("manger: %v\n", m)
// 	m.SelectWorker()
// 	m.UpdateTasks()
// 	m.SendWork()

// 	n := node.Node{
// 		Name:   "Node-1",
// 		Ip:     "192.168.1.1",
// 		Cores:  4,
// 		Memory: 1024,
// 		Disk:   25,
// 		Role:   "worker",
// 	}
// 	fmt.Printf("node: %v\n", n)
// }
