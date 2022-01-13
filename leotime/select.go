package leotime

import (
	"fmt"
	"time"
)

type TicketTask struct {
	ch chan bool
}

// open a timer task
// close : TicketTask.Close()
// we should add recover in 'task func()' avoid err
func StartTickerTask(d time.Duration, startFirst bool, task func()) *TicketTask {
	ticker := new(TicketTask)
	t := time.NewTicker(d)
	ch := make(chan bool)
	go runTickerTask(t, ch, task)
	if startFirst {
		go task()
	}
	ticker.ch = ch
	return ticker
}

func runTickerTask(t *time.Ticker, ch chan bool, task func()) {
	defer t.Stop()
	for {
		select {
		case _ = <-t.C:
			if task != nil {
				go task()
			}
		case stop := <-ch:
			if stop {
				fmt.Println("Ticker Stop")
				return
			}
		}
	}
}

func (t *TicketTask) Close() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("channel has Closed", err)
		}
	}()
	t.ch <- true
	close(t.ch)
}

// open a timer task
// close : TicketTask.Close()
// we should add recover in 'task func()' avoid err

// Used for tasks with short intervals but may take a long time;
// the next task will start after d time after the task is completed
func StartTimerTask(d time.Duration, startFirst bool, task func()) *TicketTask {
	ticker := new(TicketTask)
	t := time.NewTimer(d)
	ch := make(chan bool)
	go runTimerTask(t, d, ch, task)
	if startFirst {
		go task()
	}
	ticker.ch = ch
	return ticker
}

func runTimerTask(t *time.Timer, d time.Duration, ch chan bool, task func()) {
	defer t.Stop()
	for {
		select {
		case _ = <-t.C:
			if task != nil {
				go func() {
					defer func() {
						err := recover()
						if err != nil {
							t.Reset(d)
						}
					}()
					task()
					t.Reset(d)
				}()
			}
		case stop := <-ch:
			if stop {
				fmt.Println("Ticker Stop")
				return
			}
		}
	}
}
