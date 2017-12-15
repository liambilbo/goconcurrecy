package goroutines

import (
	"fmt"
	"testing"
	//"github.com/stretchr/testify/assert"
	//"sync"
	"sync"
)
func TestConcurrency(t *testing.T) {
}

func ExampleForkJoint() {
	var wg sync.WaitGroup
	sayHello:=func() {
		defer wg.Done()
		fmt.Println("Hello!")
	}
	wg.Add(1)
	go sayHello()
	wg.Wait()
	// Output:
	// Hello!
}


func ExampleClousureLoopError() {
	var wg sync.WaitGroup
	for _,salutation :=range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
			return
		}()
	}
	wg.Wait()
	// Output:
	// good day
	// good day
	// good day
}

func ExampleClousureLoop() {
	var wg sync.WaitGroup
	for _,salutation :=range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
			return
		}(salutation)
	}
	wg.Wait()
	// Output:
	// hello
	// greetings
	// good day
}




