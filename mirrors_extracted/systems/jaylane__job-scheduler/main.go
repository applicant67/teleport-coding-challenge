package main

import (
	"fmt"

	j "github.com/jaylane/job-scheduler/pkg/job"
	"github.com/jaylane/job-scheduler/pkg/worker"
	"github.com/jaylane/job-scheduler/pkg/worker/config"
)

func main() {
	w := worker.NewWorker(config.NewConfig())
	id, _ := w.StartJob(j.Command{"/bin/echo", []string{"hello", "world"}})

	fmt.Println(id)

}
