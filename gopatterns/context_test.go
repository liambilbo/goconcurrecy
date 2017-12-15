package gopatterns

import (
	"testing"
	"context"
	"time"
	"fmt"
	"sync"
)

func TestContext(t *testing.T) {

	ctx,cancel:=context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err:=printGreeting(ctx) ; err!=nil {
			fmt.Println("Error printGreeting: ",err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err:=printFarewell(ctx) ;err!=nil {
			fmt.Println("Error printFarewell: ",err)
			cancel()
		}
	}()

	wg.Wait()




}


func printFarewell(ctx context.Context) error {
	message, err:=genFarewell(ctx)

	if err!=nil {
		return err
	}
	fmt.Println(message)
	return nil

}


func printGreeting(ctx context.Context) error {
	message,err:=genGreeting(ctx)
	if err!=nil {
		return err
	}
	fmt.Println(message)
	return nil
}

func genGreeting(ctx context.Context) (string,error) {
	ctx,cancel:=context.WithTimeout(ctx,1 * time.Second)
	defer cancel()
	switch locale,err:=locale(ctx); {
	case err!=nil:
		return "",err
	case locale=="EN/US":
		return "hello",nil
	}

	return "",fmt.Errorf("unssuported locale")

}


func genFarewell(ctx context.Context) (string,error){
	switch locale ,err:=locale(ctx);{
		case err!=nil:
			return "",err
		case locale=="EN/US":
			return "goodbye",nil
	}

	return "",fmt.Errorf("unssuported locale")
}


func locale(ctx context.Context) (string,error) {

	select {
		case <-ctx.Done():
			return "",ctx.Err()
		case <-time.After(1 * time.Minute):
	}

	return "EN/US",nil

}
