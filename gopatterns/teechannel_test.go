package gopatterns

import (
	"testing"
	"fmt"
)

var orDone func(done <-chan interface{},inputStream <-chan interface{}) <-chan interface{}
var tee func(done <-chan interface{},inputStream <-chan interface{}) (_,_<-chan interface{})


func init() {

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

	orDone = func(done <-chan interface{},inputStream <-chan interface{}) <-chan interface{} {
		outputStream := make(chan interface{})

		go func() {
			defer close(outputStream)

			for  {
				select {
					case <-done:
						return
					case c,ok := <-inputStream:
						if ok==false{
							return
						}
						select {
							case outputStream <-c:
							case <-done:
								return
						}
				}
			}
		}()

		return outputStream
	}

	tee = func(done <-chan interface{},inputStream <-chan interface{}) (_,_<-chan interface{}) {
		out1:=make(chan interface{})
		out2:=make(chan interface{})

		go func() {
			defer close(out1)
			defer close(out2)


			for val:=range orDone(done,inputStream){
				var out1,out2 = out1,out2
				for i:=0; i<2 ; i++ {
					select {
					case <-done:
					case out1<-val:
						out1=nil
					case out2<-val:
						out2=nil
					}
				}
			}
		}()


		return out1,out2
	}
}


func TestTee(t *testing.T){

	done:=make(chan interface{})
	defer close(done)
	out1,out2:= tee(done,take(done,repeat(done,1,2),2))

	for value :=range out1{

		fmt.Printf("Values %v-%v \n",value,<-out2)

	}

}