package runner

import (
	"os/signal"
	"os"
	"time"

	"errors"
)

var ERR_TIMEOUT = errors.New("timeout")
var ERR_INTER = errors.New("interupter")

type TaskContext struct {
	tasks          []func(delta int) int
	interuptChan   chan os.Signal
	exitChan       chan int
	timeExpireChan <-chan time.Time
}

func NewTaskContext(timeout time.Duration) *TaskContext {
	return &TaskContext{
		interuptChan:   make(chan os.Signal, 1),
		exitChan:       make(chan int),
		timeExpireChan: time.After(timeout),
	}
}

func (this *TaskContext) AddTask(tasks ... func(delta int) int) {
	this.tasks = append(this.tasks, tasks...)
}

func taskRoutine(this *TaskContext) {
	for k, v := range this.tasks {
		if this.scanSysmsg() {
			this.exitChan <- 1
		}
		v(k)
	}

}

func (this *TaskContext) scanSysmsg() bool {
	select {
	case <-this.interuptChan:
		signal.Stop(this.interuptChan)
		return true
	default:
		return false
	}
}

func (this *TaskContext) Start() error {
	signal.Notify(this.interuptChan, os.Interrupt)
	go taskRoutine(this)

	select {
	case res := <-this.exitChan:
		if res == 1 {
			return ERR_INTER
		}
		return nil
	case <-this.timeExpireChan:
		return ERR_TIMEOUT
	}
}
