package gopatterns



import (
	"testing"
)

var srepeat func(done <-chan interface{},data ... string) <-chan string
var stake func(done <-chan interface{}, inputStream <-chan string, number int) <-chan string

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


	srepeat= func(done <-chan interface{},data ... string) <-chan string {
		outputStream:=make(chan string)

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


	stake= func(done <-chan interface{}, inputStream <-chan string, number int) <-chan string {
		outputStream:=make(chan string)

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

}

func BenchmarkGeneric(b *testing.B) {
	done:=make(chan interface{})
	defer close(done)

	b.StartTimer()

	for range toString(done,take(done,repeat(done,"a"),b.N)){

	}
}

func BenchmarkType(b *testing.B) {
	done:=make(chan interface{})
	defer close(done)

	b.StartTimer()

	for range stake(done,srepeat(done,"a"),b.N){

	}
}
