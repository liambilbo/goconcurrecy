package gopatterns

import (
	"testing"
	"fmt"
	"math/rand"
)

var repeat func(done <-chan interface{},data ... interface{}) <-chan interface{}
var repeatFn func(done <-chan interface{},fn ... func() interface{}) <-chan interface{}
var take func(done <-chan interface{}, inputStream <-chan interface{}, number int) <-chan interface{}
var toString func(done <-chan interface{}, inputStream <-chan interface{}) <-chan string

func init() {

	repeat= func(done <-chan interface{},data ... interface{}) <-chan interface{} {
		outputStream:=make(chan interface{})

		go func() {
			defer close(outputStream)
			for  {
				for _,v:=range data {
					select {
					case <-done:
						return
					case outputStream <- v:
					}
				}
			}
		}()
		return outputStream
	}

	repeatFn= func(done <-chan interface{},fn ... func() interface{}) <-chan interface{} {
		outputStream:=make(chan interface{})

		go func() {
			defer close(outputStream)
			for  {
				for _,v:=range fn {
					select {
					case <-done:
						return
					case outputStream <- v():
					}
				}
			}
		}()
		return outputStream
	}

	take= func(done <-chan interface{}, inputStream <-chan interface{}, number int) <-chan interface{} {
		outputStream:=make(chan interface{})

		go func() {
			defer close(outputStream)
			for i:=0;i<number ; i++  {
				select {
				case <-done:
					return
				case outputStream <- <- inputStream:
				}
			}
		}()

		return outputStream
	}

	toString=func(done <-chan interface{}, inputStream <-chan interface{}) <-chan string {
		outputStream:=make(chan string)
		go func() {
			defer close(outputStream)
			for v := range inputStream {
				select {
					case <-done:
						return
					case outputStream <- v.(string):
				}
			}
		}()
		return outputStream
	}

}

func GeneratorNumber(t *testing.T) {

	done:=make(chan interface{})
	defer close(done)

	for v:= range take(done,repeat(done,1),10) {
		fmt.Println(v)
	}


}

func GeneratorFn(t *testing.T) {

	done:=make(chan interface{})
	defer close(done)

	randomInt:=func() interface{} {
		return rand.Int()
	}

	for v:= range take(done,repeatFn(done,randomInt),10) {
		fmt.Println(v)
	}


}


func TestGeneratorToString(t *testing.T) {

	done:=make(chan interface{})
	defer close(done)

    var message string
	for v:= range toString(done,take(done,repeat(done,"I","am."),11)) {
		message = message + v
	}

	fmt.Println(message)
}
