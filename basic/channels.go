// _Channels_ are the pipes that connect concurrent
// goroutines. You can send values into channels from one
// goroutine and receive those values into another
// goroutine. Channels are a powerful primitive that
// underly much of Go's functionality.

package main

import "fmt"
import "time"
import "runtime"
import "os"

func main() {
	fmt.Println(runtime.NumCPU())

	// Create a new channel with `make(chan val-type)`.
	// Channels are typed by the values they convey.
	sendmessage := make(chan string,10)
	total := make(chan int,10)
	quit := make(chan struct{})

	// _Send_ a value into a channel using the `channel <-`
	// syntax. Here we send `"ping"`  to the `messages`
	// channel we made above, from a new goroutine.
	var a = func(name int) {
		cont := true
		for cont {
			//fmt.Println("In a")
			select {
			case sendmessage <- "ping":
			case <-quit:
				fmt.Println("breaking first quit a name=", name)
				cont = false
			}
			if cont {
				select {
				case total <- 1:
				case <-quit:
					fmt.Println("breaking second quit name=", name)
					cont = false
				}
			}
		}
		fmt.Println("Ending a",name)
	}

	go func() {
		var atotal = 0
		cont := true
		for cont {
			go a(atotal)
			atotal++
			time.Sleep(time.Second)
			select {
			case <-quit:
				fmt.Println("Breaking in go a")
				cont = false
			default:
			}

		}
		close(total)
		fmt.Println("Ending go a")
	}()

	var b = func(name int) {
		fi, _ := os.Open("/dev/null")
		for cont:=true;cont; {
			var msg1 string
			select {
			case msg1 = <-sendmessage:
				fmt.Fprintln(fi, msg1)
			case <-quit:
				fmt.Println("Breaking in b name=", name)
				cont = false
			}
		}
		fi.Close()
		fmt.Println("Ending b name",name)
	}
	go func() {
		var btotal = 0
		for cont := true; cont; {
			go b(btotal)
			btotal++
			time.Sleep(time.Second)
			select {
			case <-quit:
				fmt.Println("Breaking in go b")
				cont = false
			default:
			}
		}
		fmt.Println("Ending go b")
	}()
	cont := true
	var t int
	go func() {
		for cont {
			select {
			case temp := <-total:
				t += temp
			case <-quit:
				fmt.Println("Breaking in counting")
				cont = false
			}
		}
		fmt.Println("Ending counting")
	}()

	// The `<-channel` syntax _receives_ a value from the
	// channel. Here we'll receive the `"ping"` message
	// we sent above and print it out.
	time.Sleep(time.Second * 10)
	close(quit)
	fmt.Println("Total Messages", t)
	time.Sleep(time.Second)
}
