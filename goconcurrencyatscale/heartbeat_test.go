package goconcurrencyatscale

import (
	"testing"
	"time"
	"fmt"
)

func TestHeartBeat(t *testing.T) {

	dowork:= func(done <- chan interface{},pulseIntervale time.Duration) (<-chan interface{},<-chan time.Time) {
		heartbeat:=make(chan interface{})
		results:=make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse:=time.Tick(pulseIntervale)
			workGen:=time.Tick(2 * pulseIntervale)


			sendPulse:=func() {
				select {
				    case heartbeat<- struct {}{}:

					default:
				}
			}

			sendResult:= func(time time.Time) {
				select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results<-time:
						return
					}
			}

			for  {

				select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case work:=<-workGen:
						sendResult(work)
				}

			}
		}()

		return heartbeat,results
	}


	done:=make(chan interface{})

	time.AfterFunc(10 * time.Second, func() {
		close(done)
	})



	const timeout=time.Duration(2 * time.Second)
	hearbeat,work:=dowork(done,timeout/2)

	for {
		select {
			case _,ok:=<-hearbeat:
				if ok==false {
					return
				}
				fmt.Println("Pulse \n")
			case v,ok:=<-work:
				if ok==false {
					return
				}
				fmt.Printf("Work %v \n",v)

		case <-time.After(2 * time.Second):
			return
			
		}

	}






}
