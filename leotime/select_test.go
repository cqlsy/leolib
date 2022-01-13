package leotime

import (
	"testing"
	"time"
)

var ticker *TicketTask

func TestSelect(t *testing.T) {
	sss()
	time.Sleep(20 * time.Second)
	ticker.Close()
	println("end1")
	time.Sleep(4 * time.Second)
	ticker.Close()
	println("end2")
	time.Sleep(4 * time.Second)
	println("end3")
}

func sss() {
	var in = 0
	ticker = StartTickerTask(2*time.Second, true, func() {
		println(DateFormat(time.Now(), "mm:ss"))
		in++
		if in > 1 {
			//panic("dsadsad")
		}
		time.Sleep(3 * time.Second)
		println("now: ", in)
	})
}


func TestTimer(t *testing.T) {
	timer()
	time.Sleep(20 * time.Second)
	ticker.Close()
	println("end1")
	time.Sleep(4 * time.Second)
	ticker.Close()
	println("end2")
	time.Sleep(4 * time.Second)
	println("end3")
}

func timer() {
	var in = 0
	ticker = StartTimerTask(2*time.Second, true, func() {
		println(DateFormat(time.Now(), "mm:ss"))
		in++
		if in > 1 {
			//panic("dsadsad")
		}
		time.Sleep(3 * time.Second)
		println("now: ", in)
	})
}