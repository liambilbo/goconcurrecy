package goroutines

import (
	"sync"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"text/tabwriter"
	"os"
	"math"
)




func TestMutex(t *testing.T) {

	var count int
	var lock sync.Mutex

	increment:=func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing %d \n",count)
	}


	decrement:=func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing %d \n",count)
	}

	var wg sync.WaitGroup

	//Increment
	for i:=0 ; i<=5 ; i++ {
		wg.Add(1)
		go func() {
			increment()
			wg.Done()
		}()
	}

	//Increment
	for i:=0 ; i<=5 ; i++ {
		wg.Add(1)
		go func() {
			decrement()
			wg.Done()
		}()

	}

	wg.Wait()

	fmt.Printf("Count equal %d \n" ,count)

	assert.EqualValues(t,0,count,"Error")


}

func TestRWMutext(t *testing.T) {

	producer:=func(wg *sync.WaitGroup, lock sync.Locker) {
		defer wg.Done()
		for i:=5 ; i > 0 ; i--{
			lock.Lock()
			lock.Unlock()
			time.Sleep(1)
		}
	}

	observer:= func(wg *sync.WaitGroup, lock sync.Locker) {
		defer wg.Done()
		lock.Lock()
		defer lock.Unlock()
	}

	test := func (count int,mutex ,rwmutext sync.Locker) time.Duration{
		var wg sync.WaitGroup
		
		wg.Add(count + 1)
		sinceTime:=time.Now()
		go producer(&wg,mutex)
		for  i:=count ; i >0 ; i-- {
			go observer(&wg,rwmutext)
		}
		wg.Wait()
		return time.Since(sinceTime)

	}



	writer := tabwriter.NewWriter(os.Stdout,0,1,2,' ',0)

	defer writer.Flush()

	fmt.Fprintf(writer,"Readers\tRWMutex\tMutex\n")

	var m sync.RWMutex

	for i:=0 ; i<20; i++ {
		count := int(math.Pow(2,float64(i)))
		fmt.Fprintf(writer,
			"%d\t%v\t%v\n",
				count,
			    test(count,&m,m.RLocker()),
			    test(count,&m,&m)	)
	}

}




