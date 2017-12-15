package goroutines

import (
	"testing"
	"sync"
	"fmt"
	"net"
	"time"
	"log"
	"io/ioutil"
)

func Pool(t *testing.T) {

	pool := sync.Pool{New: func() interface{} {
		fmt.Println("Creating a new interface")
		return struct {}{}
	}}

	pool.Get()
	instance:=pool.Get()
	pool.Put(instance)
	pool.Get()

	//Output :
	//Creating a new interface
	//Creating a new interface

}

func Get(t *testing.T) {

	numCalculators :=0

	pool:=&sync.Pool{New: func() interface{} {
		numCalculators +=1
		mem:=make([]byte,1024)
		return &mem
	} }


	pool.Put(pool.Get())
	pool.Put(pool.Get())
	pool.Put(pool.Get())
	pool.Put(pool.Get())

	const numworker = 1024 * 1024
	var wg sync.WaitGroup

	wg.Add(numworker)

	for i:=0 ; i < numworker ; i++ {
		go func() {
			defer wg.Done()
			mem:=pool.Get().(*[]byte)
			defer pool.Put(mem)
		}()
	}

	wg.Wait()

	fmt.Printf("%d calculators were created! \n" , numCalculators)

}


func TestNew(t *testing.T) {

	var numCalculators int

	pool:=&sync.Pool{New: func() interface{} {
		numCalculators +=1
		mem:=make([]byte,1024)
		return &mem
	} }


	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())

	const numworker = 1024 * 1024
	var wg sync.WaitGroup

	wg.Add(numworker)

	for i:=0 ; i < numworker ; i++ {
		go func() {
			defer wg.Done()
			mem:=pool.Get().(*[]byte)
			defer pool.Put(mem)
		}()
	}

	wg.Wait()

	fmt.Printf("%d calculators were created! \n" , numCalculators)

}


func connectToService() interface{} {
	time.Sleep(1)
	return struct {}{}
}


func warmServiceConnCache() *sync.Pool{
	pool:=&sync.Pool{New:connectToService}
	for i:=10 ; i>0 ; i-- {
		pool.Put(pool.New)
	}
	return pool
}


func init() {
	daemonStarted:=startNetWorkDaemon()
	daemonStarted.Wait()
}

func startNetWorkDaemon() *sync.WaitGroup {

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {

		server , error := net.Listen("tcp","localhost:8090")

		pool:=warmServiceConnCache()

		if error!=nil {
			log.Fatalf("Error Server %v",error.Error())
		}

		wg.Done()
		defer server.Close()

		for  {
			conn, error := server.Accept()

			if error!=nil {
				log.Printf("Error Accept %v",error.Error())
				continue
			}


			connectToService()

			connPoll:=pool.Get()
			pool.Put(connPoll)


			fmt.Fprint(conn,"")
			conn.Close()
		}
	}()

	return &wg

}

func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8090")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
