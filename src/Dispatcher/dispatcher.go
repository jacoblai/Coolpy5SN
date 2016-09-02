package Dispatcher

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(maxWorkers int, cphttp string) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		cpHttp: cphttp,
		JobQueue:   make(chan Job, 10000),
		maxWorkers: maxWorkers,
		workerPool: workerPool,
	}
}

type Dispatcher struct {
	cpHttp     string
	workerPool chan chan Job
	maxWorkers int
	JobQueue   chan Job
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i + 1, d.cpHttp, d.workerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func() {
				workerJobQueue := <-d.workerPool
				workerJobQueue <- job
			}()
		}
	}
}
