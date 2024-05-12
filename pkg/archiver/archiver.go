package archiver

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	StatusWaiting = iota
	StatusRunning
	StatusComplete
)

type Archiver struct {
	status   int
	progress float64
	filepath string
	context  context.Context
	cancel context.CancelFunc
}

func NewArchiver() *Archiver {
	return &Archiver{
		status:   StatusWaiting,
		progress: 0,
		filepath: "",
	}
}

// Why use these functions instead of just exposing the variables?
//
// 1. So you can't write to the variables from outside the package
// 2. More importantly so it matches the API as described in
//    hypermedia.systems

func (a *Archiver) Status() int {
	return a.status
}

func (a *Archiver) Progress() float64 {
	return a.progress
}

func (a *Archiver) ArchiveFile() string {
	return a.filepath
}

func task(a *Archiver) {
	for {
		select {
		case <- a.context.Done():
			fmt.Println("Archiver task canceled")
			a.cancel = nil
			a.context = nil
			return
		default:
			fmt.Printf("Archiving, progress %.1f%%\n", a.progress*100)

		}

		// incr is a random float
		incr := ((rand.Float64() / 5) + 0.1) * 0.3
		a.progress += incr

		if a.progress > 1 {
			a.progress = 1
			a.status = StatusComplete
			a.filepath = "data/contacts.csv" // randomised or whatever but not important now

			fmt.Println("Archive complete")
			return

		}
	
		time.Sleep(150 * time.Millisecond)
	}
}

func (a *Archiver) Run() error {
	switch a.status {
	case StatusRunning:
		return errors.New("Task already running")
	case StatusComplete:
		return errors.New("Task already completed, call the Reset() method to run again")
	case StatusWaiting:
		ctx, cancel := context.WithCancel(context.Background())

		a.context = ctx
		a.cancel = cancel
		a.status = StatusRunning

		go task(a)
		return nil
	default:
		return errors.New("Unknown state: "+strconv.Itoa(a.status))
	}
}

func (a *Archiver) Reset() error {
	if a.status == StatusRunning {
		if a.cancel == nil {
			return errors.New("Status is running, but cancel is not a valid function reference")
		}

		fmt.Println("Cancelling archiver")
		a.cancel()
	}

	a.filepath = ""
	a.status = StatusWaiting
	
	return nil
}
