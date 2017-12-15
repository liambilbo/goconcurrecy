package gopatterns

import (
	"testing"
	"sync"
	"math/rand"
	"runtime"
	"fmt"
)


var fanin func(done <-chan interface{},channels ... <-chan interface{}) <-chan interface{}

var fanout func(done <-chan interface{},channel <-chan interface{},num int) []<-chan interface{}

var toInt func(done<-chan interface{},inputStream <- chan interface{}) <- chan  int

var primeFind func(done<-chan interface{},inputStream <- chan int)  <- chan interface{}

func init() {


	toInt = func(done<-chan interface{},inputStream <- chan interface{}) <- chan  int {
		outputStream:=make(chan int)

		go func() {
			defer close(outputStream)

			for value :=range inputStream {
				select {
				case <-done:
					return
				case outputStream <- value.(int):
				}
			}

		}()
		return outputStream
	}

	fanin = func(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
		var wg sync.WaitGroup

		multiplexChannel:=make(chan interface{})

		multiplex:=func(c <- chan interface{}) {
			defer wg.Done()
			for v:=range c {
				select {
					case <-done:
					case multiplexChannel <- v:
				}
		    }
		}

		wg.Add(len(channels))
        for _,ch:=range channels {
        	go multiplex(ch)
		}

		go func() {
			wg.Wait()
			defer close(multiplexChannel)
		}()



		return multiplexChannel
	}

	primeFind= func(done<-chan interface{},inputStream <- chan int)  <- chan interface{} {
		outputStream := make(chan interface{})

		go func() {
			defer close(outputStream)
			for integer:=range inputStream {
				prime:=false
				for divisor:=integer -1;divisor>1 ; divisor--  {
					if integer%divisor==0 {
						prime=true
						break
					}
				}
				if prime {
					select {
						case <-done:
							return
						case outputStream <- integer:

					}

				}

			}
		}()

		return outputStream

	}


}

func TestFanoutin(t *testing.T) {

	done:=make(chan interface{})
	defer close(done)

	fn:=func() interface{} { return rand.Intn(50000000)}

	randStream:=toInt(done,repeatFn(done,fn))

	numfinders:=runtime.NumCPU()

	finders:=make([]<-chan interface{},numfinders)

	for i:=0; i<numfinders;i++  {
		finders[i]=primeFind(done,randStream)
	}

	for prime:= range take(done,fanin(done,finders ...),10) {
		fmt.Printf("\t%d\n", prime)

	}











}
