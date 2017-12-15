package goroutines

import (
	"testing"
	"fmt"
	"sync"
	"bytes"
	"os"
)

func Simple(t *testing.T) {

	datach:=make(chan string)

	go func() {
		datach<-"good morning"
	}()

	gretting,ok:=<-datach
	fmt.Printf("(%v) %v)",ok,gretting)

}


func CloseChannel(t *testing.T) {
	datach:=make(chan int)
	close(datach)
	value,ok:=<-datach
	fmt.Printf("(%v) %v)",ok,value)
}


func IntChannel(t *testing.T) {
	datach:=make(chan int)

	go func() {
		for i:=1;i<=5 ;i++  {
			datach<-i
		}
		defer close(datach)
	}()

	for value := range datach {
		fmt.Printf("%v \n",value)
	}
}

func UnblockChannel(t *testing.T) {
	datach:=make(chan interface{})

	var wg sync.WaitGroup

	for i:=1;i<=5 ;i++  {
		wg.Add(1)
		go func(data int) {
			defer wg.Done()
			<-datach
			fmt.Println(data)
		}(i)
	}

	fmt.Println("Unblocking routinies")
	close(datach)
	wg.Wait()


}

func BufferedChannel(t *testing.T) {
	var buffer bytes.Buffer
	defer buffer.WriteTo(os.Stdout)
	datach := make(chan int,4)
	go func() {
		defer close(datach)
		defer fmt.Fprintf(&buffer,"Producer Done \n")
		for i:=1; i<=5; i++  {
			fmt.Fprintf(&buffer,"Sending %v \n",i)
			datach<-i
		}
	}()

	for i:=range datach{
		fmt.Fprintf(&buffer,"Receiving %v \n",i)
	}

}


func TestChanOwner(t *testing.T){
	chanowner :=  func() <-chan int{
		datach :=make(chan int,5)

		go func() {
			defer close(datach)
			for i:=1; i<=5; i++  {
				datach <- i
			}

		}()
		return datach
	}

	stream:=chanowner()

	for val := range stream {
		fmt.Printf("Receiving %v",val)
	}

	fmt.Printf("End")

}
