package gopatterns

import (
	"testing"
	"fmt"
)

func TestPipeLine(t *testing.T) {

	generator:= func(done <-chan interface{}, ints ... int) <-chan int {
		intStream:=make(chan int)
		go func() {
			defer close(intStream)

			for _,i:=range ints {
			select {
				case <-done:
					return
				case intStream <- i:
				}
			}
			}()
		return intStream
	}

	multiply:= func(done <- chan interface{},inputStream <- chan int,multiplier int) <- chan int {
		outputStream:=make(chan int)

		go func() {
			defer close(outputStream)
			for i:= range inputStream {
				select {
					case <-done:
						return
					case outputStream <- i * multiplier :
				}

			}
		}()

		return outputStream
	}


	add:= func(done <-chan interface{} , inputSream <-chan int, adder int) <- chan int {
		outputStream:=make(chan int)
		go func() {
			defer close(outputStream)
			for i := range inputSream {
				select {
					case <- done:
						return
					case outputStream <- i + adder:
				}

			}
		}()
		return outputStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream:=generator(done,1,2,3,4)

	pipeline:= multiply(done ,add(done,multiply(done,intStream,2),1),2)

	for i:= range pipeline{
		fmt.Println(i)
	}

}
