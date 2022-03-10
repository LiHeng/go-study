package concurrent

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type Job interface {
	Do() error
}

type JobChan chan Job

type WorkerChan chan JobChan

var JobQueue = make(JobChan, 10)
var WorkerPool = make(WorkerChan, MaxWorkerPoolSize)

type Worker struct {
	JobChannel JobChan
	quit       chan bool
}

type Dispatcher struct {
	Workers []*Worker
	quit    chan bool
}

const MaxWorkerPoolSize = 3
const MaxQueueSize = 5

func NewWorker() *Worker {
	return &Worker{
		JobChannel: make(JobChan, MaxQueueSize),
		quit:       make(chan bool),
	}
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		Workers: make([]*Worker, MaxWorkerPoolSize),
		quit:    make(chan bool),
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < MaxWorkerPoolSize; i++ {
		worker := NewWorker()
		d.Workers = append(d.Workers, worker)
		worker.Start()
	}

	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				jobChan := <-WorkerPool
				jobChan <- job
			}(job)
		// stop dispatcher
		case <-d.quit:
			fmt.Printf("Quit dispatcher...\n")
			// for _, worker := range d.Workers {
			// 	worker.Stop()
			// }
			return
		}
	}
}

func (d *Dispatcher) Stop() {
	d.quit <- true
}

func (w *Worker) Start() {
	go func() {
		for {
			// 将worker 注册到work pool
			WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				if err := job.Do(); err != nil {
					fmt.Printf("excute job failed with err: %v", err)
				}
			case <-w.quit:
				fmt.Printf("Quit Worker...\n")
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.quit <- true
}

type RealJob struct {
	Name string
}

func (j *RealJob) Do() error {
	fmt.Printf("job %s execute...\n", j.Name)
	return nil
}

func TestWorker(t *testing.T) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			JobQueue <- &RealJob{
				Name: "Job " + strconv.Itoa(i*3+j),
			}
		}
	}
	dispatcher := NewDispatcher()
	go dispatcher.Run()
	time.Sleep(time.Second * 2)
	dispatcher.Stop()
	fmt.Printf("after dispatcher run\n")
}
