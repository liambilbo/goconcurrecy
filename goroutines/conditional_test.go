package goroutines

import (
	"testing"
	"sync"
	"time"
	"fmt"
)

func Conditional(t *testing.T) {

	cond:=sync.NewCond(&sync.Mutex{})
	var queue = make([]interface{},0,10)

	removeQueue := func(delay time.Duration) {
		time.Sleep(delay)
		cond.L.Lock()
		queue = queue[1:]
		fmt.Println("Removing from queue")
		cond.L.Unlock()
		cond.Signal()
	}


	for i:=0 ; i < 10 ; i++ {
		cond.L.Lock()
		for len(queue)==2 {
			cond.Wait()
		}
		fmt.Println("Adding to queue")
		queue=append(queue, struct {}{})
		go removeQueue(1*time.Second)
		cond.L.Unlock()

	}


}


func TestConditionalBroadcast(t *testing.T) {

	type Button struct{
		Clicked *sync.Cond
	}


	button:=Button{Clicked:sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond,f func()) {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			f()
		}()
		wg.Wait()
	}

	var clickedregistered sync.WaitGroup

	clickedregistered.Add(3)

	subscribe(button.Clicked,func(){
		fmt.Println("maximized Window")
		clickedregistered.Done()
	})

	subscribe(button.Clicked,func(){
		fmt.Println("close Window")
		clickedregistered.Done()
	})

	subscribe(button.Clicked,func(){
		fmt.Println("mouse clicked")
		clickedregistered.Done()
	})

	button.Clicked.Broadcast()

	clickedregistered.Wait()


}