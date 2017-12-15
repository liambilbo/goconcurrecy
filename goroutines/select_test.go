package goroutines

import (
	"testing"
	"time"
	"fmt"
)

func SimpleSelect(t *testing.T) {

	start:=time.Now()
	datach:=make(chan interface{})

	go func() {
		time.Sleep(5 * time.Second)
		close(datach)
	}()

	fmt.Println("Blocking on read")
	select {
		case <-datach :
			fmt.Printf("Unblocked after %v secondes \n", time.Since(start))

	}

	fmt.Println("Finish")
}


func TwoChannelReady(t *testing.T) {

	c1:=make(chan interface{}) ; close(c1)
	c2:=make(chan interface{}) ; close(c2)


	var count1,count2 int
	for i:=0; i<1000; i++  {
		select {
			case <- c1 :
				count1++
			case <- c2 :
				count2++
		}
	}

	fmt.Printf("Count 1 %v - Count 2 %v \n" , count1 , count2)


}


func TimeOut(t *testing.T) {

	start := time.Now()
	var c <-chan int
	select {
		case <-c:
		case <-time.After(5 * time.Second) :
			fmt.Printf("Time out after %v \n" ,time.Since(start))
	}

}

func SelectDefault(t *testing.T) {
	var c1,c2 chan int
	start := time.Now()

	select {
		case <-c1:
		case <-c2:
		default:
			fmt.Printf("in default after %v \n\n" , time.Since(start))
	}

}


func TestDefaultWithLoop(t *testing.T) {
	done:=make(chan interface{})
	var count1, count2 int

	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	loop:
	for  {
		select {
			case <-done:
				break loop
			default:
				count1++
		}
		count2++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Number of cycles %v - %v \n",count1,count2)






}
