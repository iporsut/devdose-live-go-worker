package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type worker struct {
	name     string
	duration time.Duration
	done     chan struct{}
}

func SleepContext(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}

func (w *worker) Execute(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s: DONE", w.name)
			close(w.done)
			return
		default:
			log.Printf("%s: at %s", w.name, time.Now().Format(time.RFC3339))
			SleepContext(ctx, w.duration)
		}
	}
}

func (w *worker) Wait() {
	<-w.done
}

func NewWorker(name string, d time.Duration) *worker {
	return &worker{
		done:     make(chan struct{}),
		name:     name,
		duration: d,
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	w := NewWorker("worker 1", 1*time.Second)
	go w.Execute(ctx)

	w2 := NewWorker("worker 2", 3*time.Second)
	go w2.Execute(ctx)

	w.Wait()
	w2.Wait()
}
