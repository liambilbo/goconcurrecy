package gopatterns

import (
	"testing"
	"fmt"
)

var bridge func(done <-chan interface{}, chanStream <- chan <- chan interface{}) <-chan interface{}


func init() {


	bridge = func(done <-chan interface{}, chanStream <- chan <- chan interface{}) <-chan interface{} {
		outputStream:=make(chan interface{})

		go func() {
			defer close(outputStream)

			for  {

				var valStream <-chan interface{}

				select {
					case <-done:
						return
				case maybeStream,ok:=<-chanStream:
					if ok==false {
						return
					}
					valStream=maybeStream
				}

				for stream := range orDone(done,valStream){
					select {
						case outputStream <- stream:
						case <-done:
					}

				}



			}



		}()

		return outputStream



	}


}


func TestBridge(t *testing.T){

	genvals:=func() <-chan (<-chan interface{}) {
		mainoutStream:=make(chan (<-chan interface{}))
		go func() {
			defer close(mainoutStream)
			for i:=0;i<10 ;i++  {
				genStream:=make(chan interface{},1)
				genStream <- i
				close(genStream)
				mainoutStream <- genStream
			}

		}()

		return mainoutStream
	}

	for v:=range bridge(nil,genvals()) {
		fmt.Printf("%v ",v)

	}
}
