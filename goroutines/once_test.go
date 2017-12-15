package goroutines

import (
	"testing"
	"sync"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {

	var count int

	increment := func() {
		count ++
	}

	once:=sync.Once{}

	var wg sync.WaitGroup
	wg.Add(100)
	for i:=0 ; i<100 ; i++ {
		go func(){
			defer wg.Done()
			once.Do(increment)
		}()
	}

	wg.Wait()

	fmt.Printf("Count : %d \n",count)

	assert.EqualValues(t,1,count,"Error")


}

func TestOnceTwoAction(t *testing.T) {
	var count int

	once := sync.Once{}
	increment:=func(){
		count++
	}
	decrement:=func(){
		count--
	}


	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count : %d \n",count)

	assert.EqualValues(t,1,count,"Error")

}
