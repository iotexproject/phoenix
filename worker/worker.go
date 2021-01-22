package worker

import (
	"context"
)

type Job struct {
	ID   uint64
	Data []byte
}

type JobChannel chan Job
type JobQueue chan chan Job

type worker struct {
	ID      int           // id of the worker
	JobChan JobChannel    // a channel to receive single unit of work
	Queue   JobQueue      // shared between all workers.
	Quit    chan struct{} // a channel to quit working
}

func newWorker(ID int, JobChan JobChannel, Queue JobQueue, Quit chan struct{}) *worker {
	return &worker{
		ID:      ID,
		JobChan: JobChan,
		Queue:   Queue,
		Quit:    Quit,
	}
}

// stop closes the Quit channel on the worker.
func (wr *worker) Start(ctx context.Context, f func(Job)) {
	go func() {
		for {
			// when available, put the JobChan again on the JobPool
			// and wait to receive a job
			wr.Queue <- wr.JobChan
			select {
			case <-ctx.Done():
				return
			case job := <-wr.JobChan:
				// when a job is received, process
				f(job)
			case <-wr.Quit:
				// a signal on this channel means someone triggered
				// a shutdown for this worker
				close(wr.JobChan)
				return
			}
		}
	}()
}

// stop closes the Quit channel on the worker.
func (wr *worker) Stop() {
	close(wr.Quit)
}

type Worker struct {
	workers  []*worker
	WorkChan JobChannel // client submits job to this channel
	Queue    JobQueue   // this is the shared JobPool between the workers
}

func New(workerNum int) *Worker {
	return &Worker{
		workers:  make([]*worker, workerNum),
		WorkChan: make(JobChannel),
		Queue:    make(JobQueue),
	}
}

func (w *Worker) Start(ctx context.Context, f func(Job)) *Worker {
	l := len(w.workers)
	for i := 1; i <= l; i++ {
		wrk := newWorker(i, make(JobChannel), w.Queue, make(chan struct{}))
		wrk.Start(ctx, f)
		w.workers = append(w.workers, wrk)
	}
	go w.process(ctx)
	return w
}

func (w *Worker) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-w.WorkChan: // listen to any submitted job on the WorkChan
			// wait for a worker to submit JobChan to Queue
			// note that this Queue is shared among all workers.
			// Whenever there is an available JobChan on Queue pull it
			jobChan := <-w.Queue

			// Once a jobChan is available, send the submitted Job on this JobChan
			jobChan <- job
		}
	}
}

func (w *Worker) Submit(job Job) {
	w.WorkChan <- job
}
