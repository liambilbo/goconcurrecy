package gopatterns

import (
	"testing"
	"fmt"
	"time"
	"math/rand"
)

func EndRoutine(t *testing.T) {

	dowork:=func(done <- chan interface{},
	             strings <- chan string) <-chan interface{} {
	             	terminated := make(chan interface{})
	             	go func() {
						defer fmt.Println("Do work exited")
						defer close(terminated)

						for  {
							select {
							case <-done:
								return
							case s := <- strings:
								fmt.Println(s)
							}
						}
					}()
					return terminated
	}

	done:=make(chan interface{})
	terminated:=dowork(done,nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling Do work routine")
		close(done)
	}()

	<-terminated
	fmt.Println("Done")

}

func TestChanProducer(t *testing.T) {


	producter:= func(done <-chan interface{}) <-chan int {
		intRandom:=make(chan int )
		go func() {
			defer fmt.Println("Ending producer")
			defer close(intRandom)
			for  {
				select {
				case <-done:
					return
				case intRandom <- rand.Int():
				}
			}
		}()

		return intRandom
	}

	done:=make(chan interface{})

	intStream:=producter(done)

	for i:=0;i<3 ; i++ {
		fmt.Printf("Random int %v \n",<-intStream)
	}

	close(done)

	time.Sleep(1 * time.Second)
}



